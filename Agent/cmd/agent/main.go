package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/AxDvl/GoLessons/agent/auxilaries"
)

const BASE_URL = "http://127.0.0.1:8081" //Пока захардкодим адрес сервера
const CONFIG_URL = BASE_URL + "/api/server-config"
const REG_URL = BASE_URL + "/api/reg-agent"
const TAKE_TASK_URL = BASE_URL + "/api/take-task?agentid=%d"
const SEND_STATUS_URL = BASE_URL + "/api/send-status"

type ConfigStruct struct {
	PlusDuration     int //Время вычисления операции "+" в секундах
	MinusDuration    int //Время вычисления операции "-" в секундах
	MulDuration      int //Время вычисления операции "*" в секундах
	DivideDuration   int //Время вычисления операции "/" в секундах
	AgentWaitTimeout int //Время ожидания ответа от агента (вычислителя) прежде чем признать его неактивным (в секундах)
}

func getConfig(client *http.Client) ConfigStruct {
	res := ConfigStruct{PlusDuration: 10, MinusDuration: 10, MulDuration: 10, DivideDuration: 10, AgentWaitTimeout: 60}
	req, err := http.NewRequest(http.MethodGet, CONFIG_URL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return res
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return res
	}
	defer resp.Body.Close()

	err = auxilaries.GetBodyAsJson(resp.Body, &res)
	if err != nil {
		fmt.Println(err.Error())
		//text, _ := auxilaries.GetStringFromBody(resp.Body)
		//fmt.Println(text)
	}
	return res
}

type AgentInfo struct {
	Id int
}

func regAgent(client *http.Client) (int, error) {
	req, err := http.NewRequest(http.MethodPost, REG_URL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var a AgentInfo

	err = auxilaries.GetBodyAsJson(resp.Body, &a)
	if err != nil {
		return 0, err
	}
	return a.Id, nil
}

type ExpressionInfo struct {
	Id           string
	Status       int
	ExecutorId   int
	Result       float32
	LastUpdate   time.Time
	LeftOperand  float32
	RightOperand float32
	Operator     int
}

func takeTask(client *http.Client, agentId int) (ExpressionInfo, bool, error) {
	var res ExpressionInfo
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(TAKE_TASK_URL, agentId), nil)
	if err != nil {
		return res, false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return res, false, err
	}
	defer resp.Body.Close()

	err = auxilaries.GetBodyAsJson(resp.Body, &res)
	if err != nil {
		return res, false, err
	}
	return res, true, nil
}

const (
	OPPlus int = iota
	OPMul
	OPDivide
)

func resolveTask(task ExpressionInfo, config ConfigStruct) float32 {
	var res float32
	switch task.Operator {
	case OPPlus:
		res = task.LeftOperand + task.RightOperand
		time.Sleep(time.Duration(config.PlusDuration) * time.Second)
	case OPMul:
		res = task.LeftOperand * task.RightOperand
		time.Sleep(time.Duration(config.MulDuration) * time.Second)
	case OPDivide:
		res = task.LeftOperand / task.RightOperand
		time.Sleep(time.Duration(config.DivideDuration) * time.Second)
	}
	return res
}

type AgentResult struct {
	ExpressionId string  `json:"expressionid"`
	Done         bool    `json:"done,omitempty"`
	Result       float32 `json:"result,omitempty"`
}

func SendStatus(client *http.Client, taskId string, done bool, result float32) {
	res := AgentResult{ExpressionId: taskId, Done: done, Result: result}
	data, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reqBody := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, SEND_STATUS_URL, reqBody)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	client := &http.Client{}
	agentId, err := regAgent(client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	config := getConfig(client)

	var wg sync.WaitGroup
	wg.Add(1)

	//ctx, cancel := context.WithCancel(context.Background())
	ctx := context.Background()

	go func() {
		defer wg.Done()
		internalCtx, stopSentStatus := context.WithCancel(context.Background())
		for {
			select {
			case <-ctx.Done():
				return
			default:
				task, ok, err := takeTask(client, agentId)
				if err != nil {
					fmt.Println(err.Error())
				}
				if !ok {
					time.Sleep(100 * time.Millisecond)
					continue
				}

				go func(ctx context.Context) {
					select {
					case <-ctx.Done():
						return
					case <-time.After(time.Duration(config.AgentWaitTimeout/2) * time.Second):
						SendStatus(client, task.Id, false, 0)
					}
				}(internalCtx)

				res := resolveTask(task, config)
				stopSentStatus()
				SendStatus(client, task.Id, true, res)

			}
		}
	}()

	fmt.Println(agentId)
	fmt.Println(config)
	wg.Wait()
}

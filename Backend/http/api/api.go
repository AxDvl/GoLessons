package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AxDvl/GoLessons/backend/internal/auxilaries"
	"github.com/AxDvl/GoLessons/backend/internal/storage"
)

type TestApiT = struct {
}

func NewApiHandler(ctx context.Context) (http.Handler, error) {
	serveMux := http.NewServeMux()

	path, err := os.Executable()
	if err != nil {
		path = ""
	}
	path = filepath.Dir(path)
	fmt.Println(path)

	serveMux.Handle("/", http.FileServer(http.Dir(filepath.Join(path, "web"))))

	serveMux.HandleFunc("/api/task", setTask)
	serveMux.HandleFunc("/api/server-config", serverConfig)

	//API для агентов
	serveMux.HandleFunc("/api/reg-agent", regAgent)
	serveMux.HandleFunc("/api/send-status", sendStatus) //Используется как для отправки результата, так и в качестве пинга (чтобы дать сервру понять, что агент еще на линии)
	serveMux.HandleFunc("/api/take-task", takeTask)

	return serveMux, nil
}

func setTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		taskText, _ := auxilaries.GetStringFromBody(r.Body)
		task := storage.TaskStore.AddTask(taskText)
		err := json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if r.Method == http.MethodGet {
		var tasks []storage.TaskInfo
		storage.TaskStore.Mu.RLock()
		for _, value := range storage.TaskStore.Tasks {
			tasks = append(tasks, *value)
		}
		storage.TaskStore.Mu.RUnlock()
		err := json.NewEncoder(w).Encode(tasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type AgentResult struct {
	ExpressionId string  `json:"expressionid"`
	Done         bool    `json:"done,omitempty"`
	Result       float32 `json:"result,omitempty"`
}

func sendStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var res AgentResult
		err := auxilaries.GetBodyAsJson(r.Body, &res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		storage.ExpressionStore.SetResult(res.ExpressionId, res.Result, !res.Done)
	}
}

func takeTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		params := r.URL.Query()
		agetnIdstr := params.Get("agentid")
		agentId, err := strconv.Atoi(agetnIdstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		expr := storage.ExpressionStore.TakeExpression(agentId)
		err = json.NewEncoder(w).Encode(expr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

var AgentCount int

type AgentInfo struct {
	Id int
}

func regAgent(w http.ResponseWriter, r *http.Request) {
	AgentCount++
	info := AgentInfo{Id: AgentCount}
	err := json.NewEncoder(w).Encode(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serverConfig(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(storage.Config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

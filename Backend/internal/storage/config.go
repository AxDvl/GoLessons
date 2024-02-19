package storage

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type ConfigStruct struct {
	PlusDuration     int //Время вычисления операции "+" в секундах
	MinusDuration    int //Время вычисления операции "-" в секундах
	MulDuration      int //Время вычисления операции "*" в секундах
	DivideDuration   int //Время вычисления операции "/" в секундах
	AgentWaitTimeout int //Время ожидания ответа от агента (вычислителя) прежде чем признать его неактивным (в секундах)
}

var Config ConfigStruct

func NewConfig() ConfigStruct {
	//Config будем грузить из файла config.json, если его нет или возникли какие-либо ошибки при его чтении, то вернем конфиг по-умолчаню
	res := ConfigStruct{PlusDuration: 10, MinusDuration: 10, MulDuration: 10, DivideDuration: 10, AgentWaitTimeout: 60}
	path, err := os.Executable()
	if err != nil {
		path = ""
	}
	path = filepath.Join(filepath.Dir(path), "config.json")

	jsonFile, err := os.Open(path)
	if err != nil {
		return res
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return res
	}

	json.Unmarshal(jsonData, &res)

	return res

}

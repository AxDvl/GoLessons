package storage

type ConfigStruct struct {
	PlusDuration     int //Время вычисления операции "+" в секундах
	MinusDuration    int //Время вычисления операции "-" в секундах
	MulDuration      int //Время вычисления операции "*" в секундах
	DivideDuration   int //Время вычисления операции "/" в секундах
	AgentWaitTimeout int //Время ожидания ответа от агента (вычислителя) прежде чем признать его неактивным (в секундах)
}

var Config ConfigStruct

func NewConfig() ConfigStruct {
	return ConfigStruct{PlusDuration: 10, MinusDuration: 10, MulDuration: 10, DivideDuration: 10, AgentWaitTimeout: 60}
}

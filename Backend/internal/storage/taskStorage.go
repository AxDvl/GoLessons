package storage

import (
	"strings"
	"sync"
	"time"
)

// Состояние выражения
const TaskStatusNew int = 0   //Выражение еще не обработано
const TaskStatusDone int = 1  //Выражение обработано и вычеслен его результат
const TaskStatusError int = 2 //Ошибка парсинга или вычисления (деление на 0) выражения

type TaskInfo struct {
	ID          string    //В качестве идентификатора будем использовать CleanValue
	TaskText    string    //Выражение, которое ввел пользователь со всеми пробелами
	CleanValue  string    //Выражение очищенное от пробелов
	Result      float32   //Результат вычисления выражения
	Status      int       //Состоояние выражения
	ResolveTime time.Time //Время решения (если Status=TaskStatusDone) или время ошибки (если Status=TaskStatusError),...
	//... предполагается использовать периодической очистки списка задач
}

type TaskStoreStruct struct {
	Mu    sync.RWMutex
	Tasks map[string]*TaskInfo
}

func NewStore() *TaskStoreStruct {
	return &TaskStoreStruct{Tasks: make(map[string]*TaskInfo)}
}

var TaskStore *TaskStoreStruct

// Добавляет выражение в хранилище, а если такое выражение уже есть, то возвращает его
func (store *TaskStoreStruct) AddTask(taskText string) TaskInfo {
	taskID := strings.ReplaceAll(taskText, " ", "")
	store.Mu.RLock()
	if tsk, ok := store.Tasks[taskID]; ok {
		return *tsk
	}
	store.Mu.RUnlock()

	tsk := TaskInfo{ID: taskID, TaskText: taskText, CleanValue: taskID, Result: 0, Status: TaskStatusNew}
	store.Mu.Lock()
	store.Tasks[taskID] = &tsk
	store.Mu.Unlock()
	return tsk
}

func (store *TaskStoreStruct) SetTaskWrongParseStatus(taskID string) {
	if tsk, ok := store.Tasks[taskID]; ok {
		tsk.Status = TaskStatusError
	}
}

func (store *TaskStoreStruct) SetTaskResult(taskID string, result float32) {
	if tsk, ok := store.Tasks[taskID]; ok {
		tsk.Result = result
		tsk.Status = TaskStatusDone
	}
}

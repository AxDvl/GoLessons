package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/AxDvl/GoLessons/backend/internal/auxilaries"
)

const (
	ExprStatusNew = iota
	ExprStatusProcessing
	ExprStatusDone
)

type ExpressionInfo struct {
	Status       int
	ExecutorId   int
	Result       float32
	LastUpdate   time.Time
	LeftOperand  float32
	RightOperand float32
	Operator     int
}

type ExpressionStoreStruct struct {
	Mu          sync.RWMutex
	Expressions map[string]*ExpressionInfo
}

func NewExpressionStore() *ExpressionStoreStruct {
	return &ExpressionStoreStruct{Expressions: make(map[string]*ExpressionInfo)}
}

var ExpressionStore *ExpressionStoreStruct

func (store *ExpressionStoreStruct) AddExpression(expr auxilaries.ExpressionToken) (ExpressionInfo, error) {
	id := expr.GetID()
	store.Mu.RLock()
	info, ok := store.Expressions[id]
	store.Mu.RUnlock()
	if ok {
		return *info, nil
	}

	if !expr.CanBeResolved() {
		return ExpressionInfo{}, errors.New("Выражение пока не может быть вычислено, так как не вычислены его подвыражения")
	}

	leftValue := expr.LeftOperand().Value()
	rightValue := expr.RightOperand().Value()
	operator := expr.Operator()
	if operator == auxilaries.OPMul && !expr.IsInverse() {
		if expr.LeftOperand().IsInverse() {
			leftValue = rightValue
			rightValue = expr.RightOperand().Value()
			operator = auxilaries.OPDivide
		}

		if expr.RightOperand().IsInverse() {
			operator = auxilaries.OPDivide
		}
	}

	info = &ExpressionInfo{
		Status:       ExprStatusNew,
		ExecutorId:   0,
		Result:       0,
		LastUpdate:   time.Now(),
		LeftOperand:  leftValue,
		RightOperand: rightValue,
		Operator:     operator,
	}

	store.Mu.Lock()
	store.Expressions[id] = info
	store.Mu.Unlock()
	return *info, nil
}

func (store *ExpressionStoreStruct) TakeExpression(executorId int) *ExpressionInfo {
	store.Mu.Lock()
	defer store.Mu.Unlock()
	for _, expr := range store.Expressions {
		if expr.Status == ExprStatusNew || expr.Status == ExprStatusProcessing && (time.Since(expr.LastUpdate) > time.Duration(Config.AgentWaitTimeout)*time.Second) {
			expr.Status = ExprStatusProcessing
			expr.LastUpdate = time.Now()
			return expr
		}
	}
	return nil
}

func (store *ExpressionStoreStruct) SetResult(expressionId string, result float32, timeUpdateOnly bool) {
	store.Mu.RLock()
	expr, ok := store.Expressions[expressionId]
	store.Mu.RUnlock()
	if ok && expr.Status != ExprStatusDone {
		expr.LastUpdate = time.Now()
		if !timeUpdateOnly {
			expr.Result = result
			expr.Status = ExprStatusDone
		}
	}
}

func (store *ExpressionStoreStruct) GetResult(expressionId string) (float32, bool) {
	store.Mu.RLock()
	expr, ok := store.Expressions[expressionId]
	store.Mu.RUnlock()
	if ok && expr.Status == ExprStatusDone {
		return expr.Result, true
	}
	return 0, false
}

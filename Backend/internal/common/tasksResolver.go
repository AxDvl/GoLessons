package common

import (
	"context"
	"fmt"
	"time"

	"github.com/AxDvl/GoLessons/backend/internal/auxilaries"
	"github.com/AxDvl/GoLessons/backend/internal/storage"
)

func StartResolve(ctx context.Context, store *storage.TaskStoreStruct, exprStore *storage.ExpressionStoreStruct) {
	go func(ctx context.Context, store *storage.TaskStoreStruct, exprStore *storage.ExpressionStoreStruct) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				store.Mu.RLock()
				for _, task := range store.Tasks {
					if task.Status == storage.TaskStatusNew {
						task.Status = storage.TaskStatusProcessing
						expr, err := auxilaries.BuildGraph(task.CleanValue)
						if err != nil {
							fmt.Println(err.Error())
							store.SetTaskWrongParseStatus(task.ID)
						} else {
							fmt.Println("==================")
							go func(task *storage.TaskInfo, ctx context.Context, exprStore *storage.ExpressionStoreStruct) {
								for {
									select {
									case <-ctx.Done():
										return
									default:
										value, ok := ResolveExpression(expr, exprStore)
										auxilaries.PrintToken(expr)
										if ok {
											store.SetTaskResult(task.ID, value)
											return
										}
									}
								}
							}(task, ctx, exprStore)

						}

					}
				}
				store.Mu.RUnlock()
				time.Sleep(time.Millisecond * 100)
			}

		}

	}(ctx, store, exprStore)
}

func ResolveExpression(exprToken auxilaries.Token, store *storage.ExpressionStoreStruct) (float32, bool) {
	if exprToken == nil {
		return 0, false
	}
	fmt.Println("=")

	var expr auxilaries.ExpressionToken
	var ok bool
	if expr, ok = exprToken.(auxilaries.ExpressionToken); !ok {
		if val, ok := exprToken.(auxilaries.ExpressionToken); ok {
			return val.Value(), true
		}
		return 0, false
	}

	if expr.CanBeResolved() {
		res, ok := store.GetResult(expr.GetID())
		fmt.Println(res)
		if !ok {
			store.AddExpression(expr)
		}
		return res, ok
	}

	if val, ok := ResolveExpression(expr.LeftOperand(), store); ok {
		expr.SetLeftOperand(auxilaries.NewValueToken(val, false))
	}

	if val, ok := ResolveExpression(expr.RightOperand(), store); ok {
		expr.SetRightOperand(auxilaries.NewValueToken(val, false))
	}

	return 0, false
}

var agentId int

// Пока это эмуляция агента
func StartAgent(ctx context.Context, exprStore *storage.ExpressionStoreStruct) {
	agentId++
	go func(ctx context.Context, exprStore *storage.ExpressionStoreStruct, agentId int) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				expr := exprStore.TakeExpression(agentId)
				if expr != nil {
					switch expr.Operator {
					case auxilaries.OPPlus:
						expr.Result = expr.LeftOperand + expr.RightOperand
					case auxilaries.OPMul:
						expr.Result = expr.LeftOperand * expr.RightOperand
					case auxilaries.OPDivide:
						expr.Result = expr.LeftOperand / expr.RightOperand
					}
					expr.Status = storage.ExprStatusDone
				}
				time.Sleep(100 * time.Millisecond)

			}

		}

	}(ctx, exprStore, agentId)
}

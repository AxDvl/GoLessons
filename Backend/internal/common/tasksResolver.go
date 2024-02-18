package common

import (
	"context"
	"fmt"
	"time"

	"github.com/AxDvl/GoLessons/backend/internal/auxilaries"
	"github.com/AxDvl/GoLessons/backend/internal/storage"
)

func StartResolve(ctx context.Context, store *storage.TaskStoreStruct) {
	go func(ctx context.Context, store *storage.TaskStoreStruct) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				store.Mu.RLock()
				for _, task := range store.Tasks {
					if task.Status == storage.TaskStatusNew {
						_, err := auxilaries.BuildGraph(task.CleanValue)
						if err != nil {
							fmt.Println(err.Error())
							store.SetTaskWrongParseStatus(task.ID)
						}

					}
				}
				store.Mu.RUnlock()
				time.Sleep(time.Millisecond * 100)
			}

		}

	}(ctx, store)
}

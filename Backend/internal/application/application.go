package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/AxDvl/GoLessons/backend/http/api"
	"github.com/AxDvl/GoLessons/backend/internal/common"
	"github.com/AxDvl/GoLessons/backend/internal/storage"
)

type Application struct{}

func New() *Application {
	return &Application{}
}

func (a *Application) Run(ctx context.Context) int {
	ctx, cancel := context.WithCancel(ctx)
	storage.Config = storage.NewConfig()
	storage.TaskStore = storage.NewStore()
	storage.ExpressionStore = storage.NewExpressionStore()
	common.StartResolve(ctx, storage.TaskStore, storage.ExpressionStore)
	common.StartAgent(ctx, storage.ExpressionStore) //TODO: Убрать когда появятся агенты

	handler, err := api.NewApiHandler(ctx)
	if err != nil {
		cancel()
		return 1
	}

	srv := &http.Server{Addr: ":8081", Handler: handler}
	go func() {
		srv.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	notifyCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-c
	cancel()
	//  завершим работу сервера
	srv.Shutdown(notifyCtx)
	fmt.Println("STOP")

	return 0

}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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

	return serveMux, nil
}

func setTask(w http.ResponseWriter, r *http.Request) {
	taskText, _ := auxilaries.GetStringFromBody(r.Body)
	task := storage.TaskStore.AddTask(taskText)
	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//go func(task storage.TaskInfo) {
	//	for _, r := range task.CleanValue {

	//	}

	//}()
}
package main

import (
	"context"

	"github.com/AxDvl/GoLessons/backend/internal/application"
)

func main() {
	app := application.New()
	app.Run(context.Background())
}

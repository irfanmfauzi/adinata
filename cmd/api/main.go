package main

import (
	"adinata/internal/handler"
	"fmt"
	"log/slog"
)

func main() {

	server := handler.NewServer()

	slog.Info("Starting Adinata")

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

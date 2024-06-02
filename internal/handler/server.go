package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"adinata/internal/service"
	"adinata/middleware"
)

type Server struct {
	port int

	service service.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		service: service.New(),
	}

	slogHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      middleware.LoggerMiddleware(NewServer.RegisterRoutes()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

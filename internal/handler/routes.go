package handler

import (
	"adinata/internal/handler/auth"
	"adinata/internal/handler/post"
	tagHandler "adinata/internal/handler/tag"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	authHandler.RegisterAuthRoute(mux, s.service.UserService)
	postHandler.RegisterPostHandler(mux, s.service.PosService)
	tagHandler.RegisterTagRoute(mux, s.service.TagService)

	return mux
}

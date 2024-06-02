package authHandler

import (
	"adinata/internal/model/request"
	"adinata/internal/model/response"
	authService "adinata/internal/service/auth"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type authHandler struct {
	authService authService.AuthServiceProvider
}

func RegisterAuthRoute(mux *http.ServeMux, authService authService.AuthServiceProvider) {
	handler := &authHandler{
		authService: authService,
	}

	mux.HandleFunc("POST /api/auth/login", handler.LoginHandler)
	mux.HandleFunc("POST /api/auth/register", handler.RegisterHandler)
}

func (a *authHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application-json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Fail when reading body", "Error", err)
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}
	defer r.Body.Close()
	req := request.LoginRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	ctx := r.Context()

	tokenString, code, err := a.authService.Login(ctx, req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(code)
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(
		response.LoginResponse{
			GenericResponse: response.GenericResponse{Success: true, Message: "Login Success"},
			Data:            response.TokenResponse{Token: tokenString},
		},
	)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func (a *authHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()
	req := request.RegisterRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	ctx := r.Context()
	code, err := a.authService.Register(ctx, req)
	if err != nil {
		resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: err.Error()})
		w.Write(resp)
		w.WriteHeader(code)
		return
	}

	resp, _ := json.Marshal(response.GenericResponse{Success: true, Message: "Success Register"})
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

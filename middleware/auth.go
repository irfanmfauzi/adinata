package middleware

import (
	"adinata/internal/model/response"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: "Unauthorized"})
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(resp)
			return
		}

		tokenString = strings.Split(tokenString, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("REPLACE_THIS_WITH_ENV_VAR_SECRET"), nil
		})
		if err != nil {
			slog.Error("ERROR", "error", err)
			resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: "Invalid Token"})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(resp)
			return
		}

		if !token.Valid {
			resp, _ := json.Marshal(response.GenericResponse{Success: false, Message: "Unauthorized"})
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(resp)
			return
		} else {
			claim := token.Claims.(jwt.MapClaims)
			user := claim["user"].(map[string]interface{})
			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}

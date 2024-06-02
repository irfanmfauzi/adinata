package authService

import (
	"adinata/internal/model"
	"adinata/internal/model/entity"
	"adinata/internal/model/request"
	"adinata/internal/repository"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type AuthServiceProvider interface {
	Login(ctx context.Context, request request.LoginRequest) (string, int, error)

	Register(ctx context.Context, request request.RegisterRequest) (int, error)
}

type AuthServiceConfig struct {
	Db       *sqlx.DB
	UserRepo repository.UserRepoProvider
}

type authService struct {
	db       *sqlx.DB
	userRepo repository.UserRepoProvider
}

func NewAuthService(cfg AuthServiceConfig) authService {
	return authService{
		db:       cfg.Db,
		userRepo: cfg.UserRepo,
	}
}

func (a *authService) Login(ctx context.Context, request request.LoginRequest) (string, int, error) {
	user, err := a.userRepo.FindUserByEmail(ctx, request.Email)
	if err != nil {
		msg := ""
		code := http.StatusInternalServerError

		if err == sql.ErrNoRows {
			msg = "Email or Password is Wrong"
			code = http.StatusBadRequest
			slog.Info("Failed to Find User By Email", "Error", err)
		} else {
			msg = "Something Wrong With System"
			slog.Error("Failed to Find User By Email", "Error", err)
		}

		return "", code, errors.New(msg)
	}

	if user.Password != request.Password {
		return "", http.StatusBadRequest, errors.New("Email or Password is Wrong")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaim{
		User: entity.User{
			Id:    user.Id,
			Email: user.Email,
			Role:  user.Role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24).UTC()),
		},
	})

	tokenString, err := token.SignedString([]byte("REPLACE_THIS_WITH_ENV_VAR_SECRET"))
	if err != nil {
		slog.Error("Failed to Signed String", "Error", err)
		return "", http.StatusInternalServerError, err
	}

	return tokenString, http.StatusOK, nil
}

func (a *authService) Register(ctx context.Context, request request.RegisterRequest) (int, error) {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("Failed to Begin Transaction", "Error", err)
		return http.StatusInternalServerError, err
	}
	defer tx.Rollback()

	_, err = a.userRepo.InsertUser(ctx, tx, request)
	if err != nil {
		slog.Error("Failed to Insert User", "Error", err)
		return http.StatusInternalServerError, err
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to Commit", "Error", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

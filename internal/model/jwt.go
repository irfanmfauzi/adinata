package model

import (
	"adinata/internal/model/entity"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaim struct {
	User entity.User `json:"user"`
	jwt.RegisteredClaims
}

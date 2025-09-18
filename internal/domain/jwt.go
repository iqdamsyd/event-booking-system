package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

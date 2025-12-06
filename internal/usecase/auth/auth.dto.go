package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GenerateToken struct {
	Id   int64
	Role uint
}

type Claims struct {
	Id    int64
	Email string
	Role  int
	jwt.RegisteredClaims
}

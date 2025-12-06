package auth

import (
	"context"
	"main/internal/entity"
)

type Auth interface {
	GenerateToken(ctx context.Context, data GenerateToken) (string, error)
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateResetToken(n int) (string, error)
}

type Repository interface {
	GetByLogin(ctx context.Context, login string) (entity.User, error)
	GetById(ctx context.Context, id int64) (entity.User, error)
}

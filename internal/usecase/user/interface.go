package user

import (
	"context"
	"main/internal/entity"

	"mime/multipart"
)

type Repository interface {
	Create(ctx context.Context, data Create) error
	GetByLogin(ctx context.Context, login string) (entity.User, error)
	GetById(ctx context.Context, id int64) (entity.User, error)
	Delete(ctx context.Context, id, deletedBy int64) error
	GetList(ctx context.Context, filter entity.Filter) ([]UserPreview, int, error)
	Update(ctx context.Context, id int64, data Update) (int64, error)
}

type Auth interface {
	HashPassword(password string) (string, error)
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

type File interface {
	Upload(ctx context.Context, file *multipart.FileHeader, folder string) (entity.File, error)
	Delete(ctx context.Context, url string) error
}

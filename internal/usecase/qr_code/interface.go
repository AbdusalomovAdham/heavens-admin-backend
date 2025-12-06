package qrcode

import (
	"context"
	"main/internal/entity"
	"mime/multipart"
)

type Repository interface {
	Create(ctx context.Context, roomId, userId int64, path string) (int64, string, error)
}

type File interface {
	Upload(ctx context.Context, file *multipart.FileHeader, folder string) (entity.File, error)
	Delete(ctx context.Context, url string) error
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

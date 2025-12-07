package room

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, roomData Create, userId int64) (int64, error)
	Delete(ctx context.Context, id, userId int64) error
	GetList(ctx context.Context, filter *entity.Filter) ([]RoomPreview, uint32, error)
	GetById(ctx context.Context, id int64) (RoomPreview, error)
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

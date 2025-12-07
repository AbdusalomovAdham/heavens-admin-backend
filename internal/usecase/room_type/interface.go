package roomtype

import (
	"context"
	"main/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data Create, userId int64) (int64, error)
	Update(ctx context.Context, data Create, userId, id int64) (int64, error)
	Delete(ctx context.Context, id, userID int64) error
	GetById(ctx context.Context, id int64) (entity.RoomType, error)
	GetList(ctx context.Context, filter *entity.Filter, userId int64) ([]entity.RoomType, uint32, error)
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

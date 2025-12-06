package departament

import (
	"context"
	"main/internal/entity"
	"main/internal/usecase/user"
)

type Repository interface {
	Create(ctx context.Context, departament Create, userId int64) (int64, error)
	Delete(ctx context.Context, id int64, userId int64) error
	GetById(ctx context.Context, id int64) (DepartamentById, error)
	GetList(ctx context.Context, filter user.Filter, lang string) ([]Get, int, error)
	Update(ctx context.Context, id int64, data Update, userId int64) error
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

package attendancetemplate

import (
	"context"
	"main/internal/entity"
	"time"
)

type Repository interface {
	Create(ctx context.Context, data Create, startAtTime, finishAtTime *time.Time) (int64, error)
	GetList(ctx context.Context, filter entity.Filter, lang string) ([]Get, int, error)
	GetById(ctx context.Context, id int64) (AttendanceTemplateById, error)
	Delete(ctx context.Context, id int64, userId int64) error
	Update(ctx context.Context, id int64, data Update, userId int64, startAtTime, finishAtTime *time.Time) error
}

type Auth interface {
	IsValidToken(ctx context.Context, tokenStr string) (entity.User, error)
}

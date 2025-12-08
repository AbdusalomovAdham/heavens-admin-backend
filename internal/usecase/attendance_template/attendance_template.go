package attendancetemplate

import (
	"context"
	"fmt"
	"main/internal/entity"
	"time"
)

type UseCase struct {
	repo Repository
	auth Auth
}

func NewUseCase(repo Repository, auth Auth) *UseCase {
	return &UseCase{
		repo: repo,
		auth: auth,
	}
}

func (uc *UseCase) Create(ctx context.Context, data Create, authHeader string) (int64, error) {
	var startAtTime *time.Time
	var finishAtTime *time.Time

	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	if data.StartAt != nil && *data.StartAt != "" {
		parsed, _ := time.Parse("15:04", *data.StartAt)
		now := time.Now()
		final := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)

		startAtTime = &final
	}

	if data.FinishAt != nil && *data.FinishAt != "" {
		parsed, err := time.Parse("15:04", *data.FinishAt)
		if err != nil {
			return 0, fmt.Errorf("invalid finish_at format: %v", err)
		}

		now := time.Now()
		t := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)
		finishAtTime = &t
	}

	data.CreatedBy = token.Id

	return uc.repo.Create(ctx, data, startAtTime, finishAtTime)
}

func (uc *UseCase) Update(ctx context.Context, id int64, data Update, authHeader string) error {
	var startAtTime *time.Time
	var finishAtTime *time.Time

	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	if data.StartAt != nil && *data.StartAt != "" {
		parsed, _ := time.Parse("15:04", *data.StartAt)
		now := time.Now()
		final := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)

		startAtTime = &final
	}

	if data.FinishAt != nil && *data.FinishAt != "" {
		parsed, err := time.Parse("15:04", *data.FinishAt)
		if err != nil {
			return fmt.Errorf("invalid finish_at format: %v", err)
		}

		now := time.Now()
		t := time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)
		finishAtTime = &t
	}

	return uc.repo.Update(ctx, id, data, token.Id, startAtTime, finishAtTime)
}

func (uc *UseCase) GetList(ctx context.Context, filter entity.Filter, lang string) ([]Get, int, error) {
	return uc.repo.GetList(ctx, filter, lang)
}

func (uc *UseCase) GetById(ctx context.Context, id int64) (AttendanceTemplateById, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *UseCase) Delete(ctx context.Context, id int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	return uc.repo.Delete(ctx, id, token.Id)
}

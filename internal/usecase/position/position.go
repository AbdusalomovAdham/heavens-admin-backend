package position

import (
	"context"
	"main/internal/entity"
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

func (uc *UseCase) Create(ctx context.Context, positionStatus Create, authHeader string) (int64, error) {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	return uc.repo.Create(ctx, positionStatus, token.Id)
}

func (uc *UseCase) Delete(ctx context.Context, positionId int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, positionId, token.Id)
}

func (uc *UseCase) GetById(ctx context.Context, id int64) (PositionById, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *UseCase) GetList(ctx context.Context, filter entity.Filter, lang string) ([]Get, int, error) {
	return uc.repo.GetList(ctx, filter, lang)
}

func (uc *UseCase) Update(ctx context.Context, id int64, data Update, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return uc.repo.Update(ctx, id, data, token.Id)
}

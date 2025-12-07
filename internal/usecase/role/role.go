package role

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

func (uc *UseCase) Create(ctx context.Context, role Create, authHeader string) (int64, error) {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	return uc.repo.Create(ctx, role, token.Id)
}

func (uc *UseCase) Delete(ctx context.Context, roleId int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, roleId, token.Id)
}

func (uc *UseCase) GetById(ctx context.Context, id int64) (RoleById, error) {
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

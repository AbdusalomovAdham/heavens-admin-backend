package departament

import (
	"context"
	"main/internal/usecase/user"
)

type UseCase struct {
	repo Repository
	auth Auth
}

func NewUseCase(repo Repository, auth Auth) *UseCase {
	return &UseCase{repo: repo, auth: auth}
}

func (uc *UseCase) Create(ctx context.Context, departament Create, authHeader string) (int64, error) {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	return uc.repo.Create(ctx, departament, token.Id)
}

func (uc *UseCase) Delete(ctx context.Context, departamentID int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, departamentID, token.Id)
}

func (uc *UseCase) GetById(ctx context.Context, id int64) (DepartamentById, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *UseCase) GetList(ctx context.Context, filter user.Filter, lang string) ([]Get, int, error) {
	return uc.repo.GetList(ctx, filter, lang)
}

func (uc *UseCase) Update(ctx context.Context, id int64, data Update, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return uc.repo.Update(ctx, id, data, token.Id)
}

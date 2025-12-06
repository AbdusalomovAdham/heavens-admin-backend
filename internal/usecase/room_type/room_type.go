package roomtype

import (
	"context"
	"main/internal/entity"
	"main/internal/usecase/user"
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

func (uc *UseCase) Create(ctx context.Context, roomType Create, authHeader string) (int64, error) {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}
	return uc.repo.Create(ctx, roomType, token.Id)
}

func (uc *UseCase) Update(ctx context.Context, roomType Create, id int64, authHeader string) (int64, error) {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}
	return uc.repo.Update(ctx, roomType, token.Id, id)
}

func (uc *UseCase) Delete(ctx context.Context, id int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id, token.Id)
}

func (uc *UseCase) GetById(ctx context.Context, id int64) (entity.RoomType, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *UseCase) GetList(ctx context.Context, filter *user.Filter, authHeader string) ([]entity.RoomType, uint32, error) {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return nil, 0, err
	}
	return uc.repo.GetList(ctx, filter, token.Id)
}

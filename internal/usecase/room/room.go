package room

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

func (uc UseCase) CreateRoom(ctx context.Context, room Create, authHeader string) (int64, error) {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	id, err := uc.repo.Create(ctx, room, token.Id)
	return id, err
}

func (uc UseCase) DeleteRoom(ctx context.Context, id int64, authHeader string) error {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	err = uc.repo.Delete(ctx, id, token.Id)
	return err
}

func (uc UseCase) GetList(ctx context.Context, filter *entity.Filter) ([]RoomPreview, uint32, error) {
	list, count, err := uc.repo.GetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return list, count, err
}

func (uc UseCase) GetById(ctx context.Context, id int64) (RoomPreview, error) {
	room, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return RoomPreview{}, err
	}
	return room, err
}

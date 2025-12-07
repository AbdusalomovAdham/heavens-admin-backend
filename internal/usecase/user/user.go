package user

import (
	"context"
	"main/internal/entity"
	"mime/multipart"
)

type UseCase struct {
	repo Repository
	auth Auth
	file File
}

func NewUseCase(repo Repository, auth Auth, file File) *UseCase {
	return &UseCase{repo: repo, auth: auth, file: file}
}

func (uc UseCase) Upload(ctx context.Context, file *multipart.FileHeader, folder string) (entity.File, error) {
	return uc.file.Upload(ctx, file, folder)
}

func (uc UseCase) AdminCreateUser(ctx context.Context, data Create, authHeader string) error {

	tokenDetail, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	hashPassword, err := uc.auth.HashPassword(data.Password)
	if err != nil {
		return err
	}

	data.Password = hashPassword
	data.CreatedBy = &tokenDetail.Id

	if err := uc.repo.Create(ctx, data); err != nil {
		return err
	}

	return nil
}

func (uc UseCase) AdminDeleteUser(ctx context.Context, id int64, authHeader string) error {
	tokenDetail, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	if err := uc.repo.Delete(ctx, id, tokenDetail.Id); err != nil {
		return err
	}

	return nil
}

func (uc UseCase) AdminGetUserDetail(ctx context.Context, id int64) (entity.User, error) {
	detail, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return entity.User{}, err
	}
	return detail, nil
}

func (uc UseCase) AdminGetUserList(ctx context.Context, filter entity.Filter) ([]UserPreview, int, error) {
	users, count, err := uc.repo.GetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (uc UseCase) AdminUpdateUser(ctx context.Context, id int64, data Update, authHeader string) (int64, error) {
	tokenDetail, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	oldUserdetail, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return 0, err
	}

	if data.Avatar != nil {
		if oldUserdetail.Avatar != "" {
			if err := uc.file.Delete(ctx, oldUserdetail.Avatar); err != nil {
				return 0, err
			}
		}
	}

	if data.CVFile != nil {
		if oldUserdetail.CVFile != nil {
			if err := uc.file.Delete(ctx, *oldUserdetail.CVFile); err != nil {
				return 0, err
			}
		}
	}

	if data.DimlomaFile != nil {
		if oldUserdetail.DimlomaFile != nil {
			if err := uc.file.Delete(ctx, *oldUserdetail.DimlomaFile); err != nil {
				return 0, err
			}
		}
	}

	data.UpdatedBy = &tokenDetail.Id
	userId, err := uc.repo.Update(ctx, id, data)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

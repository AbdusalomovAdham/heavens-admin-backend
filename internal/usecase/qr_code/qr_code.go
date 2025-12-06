package qrcode

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type UseCase struct {
	repo Repository
	file File
	auth Auth
}

func NewUseCase(repo Repository, file File, auth Auth) *UseCase {
	return &UseCase{
		repo: repo,
		file: file,
		auth: auth,
	}
}

func (uc *UseCase) GenerateQRCode(ctx context.Context, data Create, authHeader string) (string, error) {
	token, err := uc.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return "", err
	}

	randomName := uuid.New().String()
	path := fmt.Sprintf("../media/qr_codes/%s.png", randomName)

	err = qrcode.WriteFile(data.Url, qrcode.Medium, 256, path)
	if err != nil {
		return "", err
	}

	_, oldPath, err := uc.repo.Create(ctx, data.RoomId, token.Id, path)
	if err != nil {
		return "", err
	}

	if oldPath != "" {
		err = uc.file.Delete(ctx, oldPath)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

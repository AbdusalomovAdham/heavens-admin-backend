package auth

import (
	"context"
	"errors"
	"main/internal/entity"
)

type UseCase struct {
	auth Auth
	repo Repository
}

func NewUseCase(auth Auth, repo Repository) *UseCase {
	return &UseCase{auth: auth, repo: repo}
}

func (au UseCase) SignIn(ctx context.Context, data SignIn) (entity.User, string, error) {

	userDetail, err := au.repo.GetByLogin(ctx, data.Login)
	if err != nil {
		return entity.User{}, "", err
	}

	if userDetail.Password != "" {
		isValidPassword := au.auth.CheckPasswordHash(data.Password, userDetail.Password)
		if !isValidPassword {
			return entity.User{}, "", errors.New("password error")
		}
	} else {
		return entity.User{}, "", errors.New("user not found")
	}

	var generateTokenData GenerateToken
	generateTokenData.Id = userDetail.Id
	if userDetail.Role != nil {
		generateTokenData.Role = *userDetail.Role
	}

	userDetail.Password = ""

	token, err := au.auth.GenerateToken(ctx, generateTokenData)
	return userDetail, token, err
}

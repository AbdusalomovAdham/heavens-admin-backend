package auth

import (
	"context"
	"log"
	"main/internal/usecase/auth"
	auth_use_case "main/internal/usecase/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	useCase *auth.UseCase
}

func NewController(useCase *auth.UseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (ac Controller) SignIn(c *gin.Context) {
	var data auth_use_case.SignIn
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	log.Println("data", data)
	ctx := context.Background()
	detail, token, err := ac.useCase.SignIn(ctx, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"data":    detail,
		"message": "ok!",
	})
}

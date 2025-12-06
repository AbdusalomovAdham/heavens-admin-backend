package qrcode

import (
	"context"
	qrcode "main/internal/usecase/qr_code"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	useCase *qrcode.UseCase
}

func NewController(useCase *qrcode.UseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (as Controller) AdminGenerateQRCode(c *gin.Context) {
	var data qrcode.Create
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	ctx := context.Background()
	path, err := as.useCase.GenerateQRCode(ctx, data, authHeader)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"path": path, "message": "ok!"})

}

package user

import (
	"context"
	"main/internal/usecase/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	useCase *user.UseCase
}

func NewController(useCase *user.UseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (uc Controller) AdminCreateUser(c *gin.Context) {
	var data user.Create

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
		return
	}

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx := context.Background()

	file, _ := c.FormFile("avatar")
	if file != nil {
		filePath, err := uc.useCase.Upload(ctx, file, "../media/avatar")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		data.Avatar = filePath.Path
	}

	file, _ = c.FormFile("cv_file")
	if file != nil {
		filePath, err := uc.useCase.Upload(ctx, file, "../media/cv")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		data.CVFile = &filePath.Path
	}

	file, _ = c.FormFile("diploma_file")
	if file != nil {
		filePath, err := uc.useCase.Upload(ctx, file, "../media/diploma")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		data.DimlomaFile = &filePath.Path
	}

	if err := uc.useCase.AdminCreateUser(ctx, data, authHeader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

func (uc Controller) AdminDeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Id must be a number"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
		return
	}

	ctx := context.Background()
	if err := uc.useCase.AdminDeleteUser(ctx, id, authHeader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok!"})
}

func (uc Controller) AdminGetUserList(c *gin.Context) {
	var filter user.Filter
	query := c.Request.URL.Query()

	limitQ := query["limit"]
	if len(limitQ) > 0 {
		queryInt, err := strconv.Atoi(limitQ[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Limit must be number!",
			})
			return
		}

		filter.Limit = &queryInt
	}

	offsetQ := query["offset"]
	if len(offsetQ) > 0 {
		queryInt, err := strconv.Atoi(offsetQ[0])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Offset must be number!",
			})

			return
		}
		filter.Offset = &queryInt
	}

	ctx := context.Background()

	orderQ := query["order"]
	if len(orderQ) > 0 {
		filter.Order = &orderQ[0]
	}

	list, count, err := uc.useCase.AdminGetUserList(ctx, filter)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"data": map[string]any{
			"results": list,
			"count":   count,
		},
	})

}

func (uc Controller) AdminGetById(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Id must be number!",
		})
		return
	}

	ctx := context.Background()

	detail, err := uc.useCase.AdminGetUserDetail(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"data":    detail,
	})
}

func (uc Controller) AdminUpdateUser(c *gin.Context) {
	var data user.Update
	authHeader := c.GetHeader("Authorization")

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Id must be a number"})
		return
	}

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx := context.Background()

	file, _ := c.FormFile("avatar")
	if file != nil {
		filePath, err := uc.useCase.Upload(ctx, file, "../media/avatar")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		data.Avatar = &filePath.Path
	}

	file, _ = c.FormFile("cv_file")
	if file != nil {
		filePath, err := uc.useCase.Upload(ctx, file, "../media/cv")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		data.CVFile = &filePath.Path
	}

	file, _ = c.FormFile("diploma_file")
	if file != nil {
		filePath, err := uc.useCase.Upload(ctx, file, "../media/diploma")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		data.DimlomaFile = &filePath.Path
	}

	userId, err := uc.useCase.AdminUpdateUser(ctx, id, data, authHeader)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok!", "id": userId})
}

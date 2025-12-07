package workstatus

import (
	"context"
	"net/http"
	"strconv"

	"main/internal/entity"
	work_status "main/internal/usecase/work_status"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	useCase *work_status.UseCase
}

func NewController(useCase *work_status.UseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (ac Controller) AdminCreateWorkStatus(c *gin.Context) {
	var data work_status.Create

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	workStatusId, err := ac.useCase.Create(ctx, data, authHeader)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "id": workStatusId})
}

func (ac Controller) AdminDeleteWorkStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	if err := ac.useCase.Delete(ctx, id, authHeader); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!"})
}

func (ac Controller) AdminGetWorkStatusById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	workStatus, err := ac.useCase.GetById(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": workStatus})
}

func (ac Controller) AdminGetWorkStatusList(c *gin.Context) {
	var filter entity.Filter
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

	orderQ := query["order"]
	if len(orderQ) > 0 {
		filter.Order = &orderQ[0]
	}

	ctx := context.Background()
	lang := c.GetHeader("Accept-Language")

	workStatuses, total, err := ac.useCase.GetList(ctx, filter, lang)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": workStatuses, "total": total})
}

func (ac Controller) AdminUpdateWorkStatus(c *gin.Context) {
	var data work_status.Update
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	if err := ac.useCase.Update(ctx, id, data, authHeader); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!"})
}

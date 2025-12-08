package attendancetemplate

import (
	"context"
	"main/internal/entity"
	attendance "main/internal/usecase/attendance_template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	useCase *attendance.UseCase
}

func NewController(useCase *attendance.UseCase) *Controller {
	return &Controller{
		useCase: useCase,
	}
}

func (as *Controller) AdminCreateAttendanceTemplate(c *gin.Context) {
	var data attendance.Create

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
	id, err := as.useCase.Create(ctx, data, authHeader)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"id": id})
}

func (ac Controller) AdminGetAttendanceTemplateList(c *gin.Context) {
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

	attendance, total, err := ac.useCase.GetList(ctx, filter, lang)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": attendance, "total": total})
}

func (ac Controller) AdminDeleteAttendanceTemplate(c *gin.Context) {
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

func (ac Controller) AdminGetAttendanceTemplateById(c *gin.Context) {
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
	attendance, err := ac.useCase.GetById(ctx, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok!", "data": attendance})
}

func (ac Controller) AdminUpdateAttendanceTemplate(c *gin.Context) {
	var data attendance.Update
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

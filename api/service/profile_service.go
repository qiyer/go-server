package service

import (
	"go-server/domain"
	"go-server/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Fetch(c *gin.Context) {
	userID := c.GetString("x-user-id")

	profile, err := repository.GetProfileByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

package controller

import (
	"gr24/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": config.Service.FormattedNameWithVersion() + " is online!"})
}

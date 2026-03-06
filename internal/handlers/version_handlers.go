package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	APIVersion = "1.0.0"
	BuildTime  = "2026-03-06"
)

// VersionCheck godoc
// @Summary      API version info
// @Description  Returns API version and build time
// @Tags         version
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /v1/version [get]
func VersionCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":    APIVersion,
		"build_time": BuildTime,
	})
}

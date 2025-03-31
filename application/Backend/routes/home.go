package routes

import (
	"fmt"
	"net/http"

	"backend/server"

	"github.com/gin-gonic/gin"
)

func RegisterHomeRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		htmlData, err := server.LoadFrontendFile("src/html/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlData)
	})
}

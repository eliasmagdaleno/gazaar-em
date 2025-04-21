package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/gin-gonic/gin"
)

func RegisterMessagesRoutes(router *gin.Engine) {
	router.GET("/messages", func(c *gin.Context) {
		htmlData, err := core.LoadFrontendFile("src/views/messages.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlData)
	})
}
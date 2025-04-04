package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/gin-gonic/gin"
)

func RegisterVPRoutes(router *gin.Engine) {
	router.GET("/vp", func(c *gin.Context) {
		htmlData, err := core.LoadFrontendFile("src/html/vp_home.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlData)
	})
}

package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterErrorRoutes(router *gin.Engine) {
	router.GET("/error", func(c *gin.Context) {
		errorTitle := c.Query("title")
		errorMessage := c.Query("message")

		errorTemplate, err := core.LoadFrontendFile("src/views/partials/errorcard.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading error template: %v", err))
			return
		}

		content, err := raymond.Render(errorTemplate, map[string]interface{}{
			"errorTitle":   errorTitle,
			"errorMessage": errorMessage,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering error template: %v", err))
			return
		}

		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			return
		}

		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "Error",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})
}

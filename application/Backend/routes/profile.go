package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterProfileRoutes(router *gin.Engine) {
	router.GET("/profile", RandomProductMiddleware(), RandomEventMiddleware(), func(c *gin.Context) {
		username := "John Doe"
		bio := "Welcome to my profile! I love selling and hosting events."

		items, _ := c.Get("productCards")
		events, _ := c.Get("eventCards")

		profileTemplate, err := core.LoadFrontendFile("src/views/profile.hbs")
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error loading profile template: %v", err))
			return
		}

		content, err := raymond.Render(profileTemplate, map[string]interface{}{
			"username": username,
			"bio":      bio,
			"items":    items,
			"events":   events,
			"section":  c.Query("section"),
		})
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error rendering profile content: %v", err))
			return
		}

		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			return
		}

		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "Profile",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})
}

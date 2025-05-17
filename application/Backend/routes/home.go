package routes

import (
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

// Update RegisterHomeRoutes to use ProductCardMiddleware and EventCardMiddleware
func RegisterHomeRoutes(router *gin.Engine) {
	router.GET("/", RandomEventMiddleware(), RandomProductMiddleware(), func(c *gin.Context) {
		products, _ := c.Get("productCards")
		events, _ := c.Get("eventCards")

		indexTemplate, err := core.LoadFrontendFile("src/views/index.hbs")
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, "Failed to load index template")
			return
		}

		content, err := raymond.Render(indexTemplate, map[string]interface{}{
			"title":    "Home",
			"products": products,
			"events":   events,
		})
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, "Failed to render index template")
			return
		}

		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, "Failed to load layout template")
			return
		}

		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "Home",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, "Failed to render layout template")
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})
}

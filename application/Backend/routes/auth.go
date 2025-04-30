package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up the /login route
func RegisterAuthRoutes(router *gin.Engine) {

	router.POST("/logout", func(c *gin.Context) {
		fmt.Printf("This is working?\n");

		c.Redirect(http.StatusFound, "/login");
		
	})

	router.GET("/login", func(c *gin.Context) {
		// 1. Load the login template
		loginTpl, err := core.LoadFrontendFile("src/views/login.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading login template: %v", err))
			return
		}

		// 2. Render the login content (no dynamic data for now)
		content, err := raymond.Render(loginTpl, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering login template: %v", err))
			return
		}

		// 3. Load the layout
		layoutTpl, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout: %v", err))
			return
		}

		// 4. Render the final HTML
		page, err := raymond.Render(layoutTpl, map[string]interface{}{
			"title":   "Login",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, page)
	})
}

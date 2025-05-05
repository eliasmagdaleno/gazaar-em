package routes

import (
	"fmt"
	"net/http"
	"strings"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up the /login route
func RegisterAuthRoutes(router *gin.Engine) {

	router.POST("/logout", func(c *gin.Context) {
		fmt.Printf("This is working?\n")

		c.Redirect(http.StatusFound, "/login")

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

	router.POST("/login", func(c *gin.Context) {
		// Extract login credentials from the form
		email := c.PostForm("email")
		password := c.PostForm("password")

		// Validate credentials (replace with actual authentication logic)
		if email == "test@example.com" && password == "password" {
			// Redirect to the home page upon successful login
			c.Redirect(http.StatusFound, "/")
		} else {
			// Redirect back to the login page with an error message
			c.Redirect(http.StatusFound, "/login?error=invalid_credentials")
		}
	})

	router.GET("/register", func(c *gin.Context) {
		registerTpl, err := core.LoadFrontendFile("src/views/registration.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading registration template: %v", err))
			return
		}

		content, err := raymond.Render(registerTpl, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering registration template: %v", err))
			return
		}

		layoutTpl, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout: %v", err))
			return
		}

		page, err := raymond.Render(layoutTpl, map[string]interface{}{
			"title":   "Register",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, page)
	})

	// Removed unused variables 'name' and 'password' from the POST /register handler
	router.POST("/register", func(c *gin.Context) {
		email := c.PostForm("email")

		if !strings.HasSuffix(email, "@sfsu.edu") {
			c.String(http.StatusBadRequest, "Invalid email domain. Please use an @sfsu.edu email.")
			return
		}

		// Add logic to save user data to the database

		c.Redirect(http.StatusFound, "/login")
	})
}

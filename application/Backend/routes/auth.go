package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAuthRoutes sets up the /login route
func RegisterAuthRoutes(router *gin.Engine) {
	// GET /login - Show login form
	router.GET("/login", func(c *gin.Context) {
		// 1. Load the login template
		loginTpl, err := core.LoadFrontendFile("src/views/login.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading login template: %v", err))
			return
		}

		// 2. Render the login content with CSRF token
		content, err := raymond.Render(loginTpl, map[string]interface{}{
			"csrf": c.GetString("csrf_token"),
		})
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

	// POST /login - Handle login form submission
	router.POST("/login", func(c *gin.Context) {
		// 1. Get and validate input
		email := strings.TrimSpace(c.PostForm("email"))
		password := c.PostForm("password")

		// Validate email format
		if !isValidEmail(email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		// 2. Use parameterized query to prevent SQL injection
		var hashedPassword string
		err := database.DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				// Don't reveal whether the email exists or not
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// 3. Compare passwords using bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// 4. Set session cookie
		c.SetCookie("session", "authenticated", 3600, "/", "", true, true)

		// 5. Redirect to home page
		c.Redirect(http.StatusFound, "/")
	})

	router.POST("/logout", func(c *gin.Context) {
		// Clear the session cookie
		c.SetCookie("session", "", -1, "/", "", true, true)
		c.Redirect(http.StatusFound, "/login")
	})

	// Test endpoint to verify SQL injection prevention
	router.GET("/test-login", func(c *gin.Context) {
		// Test cases for SQL injection attempts
		testCases := []string{
			"test@sfsu.edu' OR '1'='1",
			"test@sfsu.edu'; DROP TABLE users; --",
			"test@sfsu.edu' UNION SELECT * FROM users; --",
		}

		results := make(map[string]string)
		for _, testCase := range testCases {
			// Try to query with the test case
			var result string
			err := database.DB.QueryRow("SELECT email FROM users WHERE email = ?", testCase).Scan(&result)

			if err == sql.ErrNoRows {
				results[testCase] = "Prevented - No results found"
			} else if err != nil {
				results[testCase] = "Prevented - Error: " + err.Error()
			} else {
				results[testCase] = "WARNING: Query succeeded!"
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "SQL Injection Prevention Test Results",
			"results": results,
		})
	})
}

// isValidEmail checks if the email is valid and is an SFSU email
func isValidEmail(email string) bool {
	// Basic email format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	// Check if it's an SFSU email
	return strings.HasSuffix(email, "@sfsu.edu")
}

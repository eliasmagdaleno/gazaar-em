package routes

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/**
 * RegisterAuthRoutes sets up the /login and /logout routes
 */
func RegisterAuthRoutes(router *gin.Engine) {
	/**
	 * GET /login - Show login form
	 */
	router.GET("/login", func(c *gin.Context) {
		showLoginPage(c, "")
	})

	/**
	 * POST /login - Handle login form submission
	 */
	router.POST("/login", func(c *gin.Context) {
		email := strings.TrimSpace(c.PostForm("email"))
		password := c.PostForm("password")

		// Validate email format
		if !isValidEmail(email) {
			showLoginPage(c, "Invalid email format.")
			return
		}

		var hashedPassword string
		err := database.DB.QueryRow("SELECT password FROM Account WHERE email_id = ?", email).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				showLoginPage(c, "Incorrect email or password.")
				return
			}
			showLoginPage(c, "Internal server error. Please try again.")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			showLoginPage(c, "Incorrect email or password.")
			return
		}
		var userID int 
		err = database.DB.QueryRow("SELECT user_id FROM Account WHERE email_id = ?", email).Scan(&userID)
		if err != nil {
			showLoginPage(c, "Internal server error. Please try again.")
			return
		}

		c.SetCookie("session", "authenticated", 7200, "/", "", true, true)
		c.SetCookie("user_id", fmt.Sprintf("%d", userID), 7200, "/", "", true, true)
		log.Println("User ID set in cookie:", userID)

		c.Redirect(http.StatusFound, "/")
	})

	/**
	 * POST /logout - Clear session cookie and redirect to login
	 */
	router.POST("/logout", func(c *gin.Context) {
		c.SetCookie("session", "", -1, "/", "", true, true)
		c.Redirect(http.StatusFound, "/login")
	})

	/**
	 * GET /test-login - Test SQL injection prevention
	 * Future note: Probably better to do tests in
	 * test file
	 */
	router.GET("/test-login", func(c *gin.Context) {
		testCases := []string{
			"test@sfsu.edu' OR '1'='1",
			"test@sfsu.edu'; DROP TABLE users; --",
			"test@sfsu.edu' UNION SELECT * FROM users; --",
		}

		results := make(map[string]string)
		for _, testCase := range testCases {
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

	/**
	* User registration makes sure each post is 
	* valid then attempts to store in data base.
	*/
	router.GET("/register", func(c *gin.Context) {
		showRegisterPage(c, "", "")
	})

	router.POST("/register", func(c *gin.Context) {
		name := strings.TrimSpace(c.PostForm("name"))
		email := strings.TrimSpace(c.PostForm("email"))
		password := c.PostForm("password")
		confirmPassword := c.PostForm("confirm-password")
	
		if name == "" {
			showRegisterPage(c, "Full name is required.", email)
			return
		}
		if !isValidEmail(email) {
			showRegisterPage(c, "A valid SFSU email is required.", email)
			return
		}
		if len(password) < 8 {
			showRegisterPage(c, "Password must be at least 8 characters.", email)
			return
		}
		if password != confirmPassword {
			showRegisterPage(c, "Passwords do not match.", email)
			return
		}
	
		var exists int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM Account WHERE email_id = ?", email).Scan(&exists)
		if err != nil {
			showRegisterPage(c, "Database error. Please try again.", email)
			return
		}
		if exists > 0 {
			showRegisterPage(c, "Email already registered.", email)
			return
		}
	
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			showRegisterPage(c, "Internal error. Please try again.", email)
			return
		}
		rand.Seed(time.Now().UnixNano())
		userID := rand.Intn(1000000000)
	
		// Insert new user
		_, err = database.DB.Exec(`
			INSERT INTO Account (user_id, user_name, email_id, password, registration_date)
			VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		`, userID, name, email, hashedPassword,)
	
		if err != nil {
			log.Println("Register DB error:", err)
			showRegisterPage(c, "Registration failed. Please try again.", email)
			return
		}
	
		c.SetCookie("session", "authenticated", 3600, "/", "", true, true)
		c.Redirect(http.StatusFound, "/")
	})
	

}

/**
 * showLoginPage renders the login page with an optional error message
 */
func showLoginPage(c *gin.Context, errorMessage string) {
	loginTpl, err := core.LoadFrontendFile("src/views/login.hbs")
	
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading login template: %v", err))
		return
	}

	content, err := raymond.Render(loginTpl, map[string]interface{}{
		"csrf":  c.GetString("csrf_token"),
		"error": errorMessage,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering login template: %v", err))
		return
	}

	layoutTpl, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout: %v", err))
		return
	}

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
}

/**
 * isValidEmail checks if the email is valid and ends with @sfsu.edu
 */
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}
	return strings.HasSuffix(email, "@sfsu.edu")
}

/**
 * showRegisterPage renders the register page with optional error and pre-filled email
 */
 func showRegisterPage(c *gin.Context, errorMessage string, email string) {
    registerTpl, err := core.LoadFrontendFile("src/views/register.hbs")
    if err != nil {
        c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading register template: %v", err))
        return
    }

    content, err := raymond.Render(registerTpl, map[string]interface{}{
        "csrf":  c.GetString("csrf_token"),
        "error": errorMessage,
        "email": email,
    })
    if err != nil {
        c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering register template: %v", err))
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
}

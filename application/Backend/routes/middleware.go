package routes

import (
	"application/Backend/core"
	"application/Backend/database"
	"log"

	"net/http"

	"github.com/aymerick/raymond"

	"github.com/gin-gonic/gin"
)

// Middleware to fetch random product cards
func RandomProductMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := "SELECT image_url, title, price FROM items WHERE LOWER(category) != 'events' ORDER BY RAND() LIMIT 20"
		rows, err := database.DB.Query(query)
		if err != nil {
			log.Printf("RandomProductMiddleware: Error executing query: %v", err)
			renderErrorPage(c, http.StatusInternalServerError, "Failed to load product cards")
			c.Abort()
			return
		}
		defer rows.Close()

		var products []map[string]interface{}
		for rows.Next() {
			var imageURL, title string
			var price float64
			if err := rows.Scan(&imageURL, &title, &price); err != nil {
				log.Printf("RandomProductMiddleware: Row scan error: %v", err)
				continue
			}
			products = append(products, map[string]interface{}{
				"thumbnail": "/assets/thumbnails/" + imageURL,
				"title":     title,
				"price":     price,
			})
		}

		c.Set("productCards", products)
		c.Next()
	}
}

// Middleware to fetch random event cards
func RandomEventMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := "SELECT image_url, title, post_date FROM items WHERE LOWER(category) = 'events' ORDER BY RAND() LIMIT 20"
		rows, err := database.DB.Query(query)
		if err != nil {
			log.Printf("RandomEventMiddleware: Error executing query: %v", err)
			c.Next()
			return
		}
		defer rows.Close()

		var events []map[string]interface{}
		for rows.Next() {
			var imageURL, title, postDate string
			if err := rows.Scan(&imageURL, &title, &postDate); err != nil {
				log.Printf("RandomEventMiddleware: Row scan error: %v", err)
				continue
			}
			if imageURL == "" {
				imageURL = "Thumbnail Unavailable"
			}
			events = append(events, map[string]interface{}{
				"thumbnail": "/assets/thumbnails/" + imageURL,
				"title":     title,
				"postDate":  postDate,
			})
		}

		c.Set("eventCards", events)
		c.Next()
	}
}

// Middleware to fetch non-random product cards
func ProductMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := "SELECT image_url, title, price FROM items WHERE LOWER(category) != 'events' LIMIT 20"
		rows, err := database.DB.Query(query)
		if err != nil {
			log.Printf("ProductMiddleware: Error executing query: %v", err)
			renderErrorPage(c, http.StatusInternalServerError, "Failed to load product cards")
			c.Abort()
			return
		}
		defer rows.Close()

		var products []map[string]interface{}
		for rows.Next() {
			var imageURL, title string
			var price float64
			if err := rows.Scan(&imageURL, &title, &price); err != nil {
				log.Printf("ProductMiddleware: Row scan error: %v", err)
				continue
			}
			products = append(products, map[string]interface{}{
				"thumbnail": "/assets/thumbnails/" + imageURL,
				"title":     title,
				"price":     price,
			})
		}

		c.Set("productCards", products)
		c.Next()
	}
}

// Helper function to render an error page
func renderErrorPage(c *gin.Context, statusCode int, message string) {
	errorTemplate, err := core.LoadFrontendFile("src/views/error.hbs")
	if err != nil {
		log.Println("Error loading error template:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	content, err := raymond.Render(errorTemplate, map[string]interface{}{
		"statusCode": statusCode,
		"message":    message,
	})
	if err != nil {
		log.Println("Error rendering error template:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(statusCode, content)
}

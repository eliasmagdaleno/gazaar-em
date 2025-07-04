package routes

import (
	"application/Backend/core"
	"application/Backend/database"
	"log"
	"strconv"

	"net/http"

	"github.com/aymerick/raymond"

	"github.com/gin-gonic/gin"
)

func UserIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userIDStr, err := c.Cookie("user_id")
        if err == nil {
            if userID, err := strconv.Atoi(userIDStr); err == nil {
                c.Set("user_id", userID)
            }
        }
        c.Next()
    }
}

// Middleware to set is_signed_in in context
func SignedInMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Ensure is_signed_in key exists in context
        session, err := c.Cookie("session")
        if err == nil && session == "authenticated" {
            c.Set("is_signed_in", true)
        } else {
            c.Set("is_signed_in", false)
        }
        c.Next()
    }
}

// Middleware to fetch random product cards
func RandomProductMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := "SELECT item_id, image_url, title, price FROM items WHERE LOWER(category) != 'events' AND approve = 1 ORDER BY RAND() LIMIT 20"
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
			var itemID int
			var imageURL, title string
			var price float64
			if err := rows.Scan(&itemID, &imageURL, &title, &price); err != nil {
				log.Printf("RandomProductMiddleware: Row scan error: %v", err)
				continue
			}
			products = append(products, map[string]interface{}{
				"item_id":   itemID,
				"thumbnail": "frontend/assets/thumbnails/" + imageURL,
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
		query := "SELECT item_id, image_url, title, post_date FROM items WHERE LOWER(category) = 'events' AND approve = 1 ORDER BY RAND() LIMIT 20"
		rows, err := database.DB.Query(query)
		if (err != nil) {
			log.Printf("RandomEventMiddleware: Error executing query: %v", err)
			c.Next()
			return
		}
		defer rows.Close()

		var events []map[string]interface{}
		for rows.Next() {
			var itemID, imageURL, title, postDate string
			if err := rows.Scan(&itemID, &imageURL, &title, &postDate); err != nil {
				log.Printf("RandomEventMiddleware: Row scan error: %v", err)
				continue
			}
			if imageURL == "" {
				imageURL = "Thumbnail Unavailable"
			}
			events = append(events, map[string]interface{}{
				"item_id":        itemID,
				"thumbnail": "frontend/assets/thumbnails/" + imageURL,
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
		query := "SELECT item_id, image_url, title, price FROM items WHERE LOWER(category) != 'events' AND approve = 1 LIMIT 20"
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
			var itemID, imageURL, title string
			var price float64
			if err := rows.Scan(&imageURL, &title, &price); err != nil {
				log.Printf("ProductMiddleware: Row scan error: %v", err)
				continue
			}
			products = append(products, map[string]interface{}{
				"item_id":   itemID,
				"thumbnail": "frontend/assets/thumbnails/" + imageURL,
				"title":     title,
				"price":     price,
			})
		}

		c.Set("productCards", products)
		c.Next()
	}
}

// Middleware to fetch a single product by ID
func ProductDetailsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		// log.Printf("ProductDetailsMiddleware: Received productID: %s", productID) // Debugging log

		query := `SELECT items.title, items.item_id, items.description, items.price, items.category, items.seller_id, items.image_url, items.address, 
		Account.user_name AS seller_name, 
			DATE_FORMAT(items.post_date, '%Y-%m-%d') AS post_date
		FROM items
		JOIN Account ON items.seller_id = Account.user_id
		WHERE item_id = ? AND approve = 1`
		// log.Printf("ProductDetailsMiddleware: Executing query: %s with productID: %s", query, productID) // Debugging log

		row := database.DB.QueryRow(query, productID)

		var product map[string]interface{}
		var title, itemID, description, category, sellerID, sellerName, imageURL, postDate, address string
		var price float64
		if err := row.Scan(&title, &itemID, &description, &price, &category, &sellerID, &imageURL, &address, &sellerName, &postDate); err != nil {
			log.Printf("ProductDetailsMiddleware: Error fetching product details: %v", err)
			renderErrorPage(c, http.StatusNotFound, "Product not found")
			c.Abort()
			return
		}

		// Map numeric address/location to human-readable string
		locationNames := []string{"Parking Garage (main) [E-F5-6]", "Student Events Center [D3-4]", " Sutro Library [F9]", "Gymnasium [G7-8]", " Creative Arts [I-J5-7]"}
		locationIdx, err := strconv.Atoi(address)
		locationName := address // fallback to raw value
		if err == nil && locationIdx >= 0 && locationIdx < len(locationNames) {
			locationName = locationNames[locationIdx]
		}

		product = map[string]interface{}{
			"title":       title,
			"id":          itemID,
			"description": description,
			"price":       price,
			"category":    category,
			"sellerID":    sellerID,
			"sellerName":  sellerName,
			"imageURL":    "frontend/assets/thumbnails/" + imageURL,
			"postDate":    postDate,
			"location":    address, // keep numeric for JS if needed
			"locationName": locationName, // human-readable for template
		}

		c.Set("productDetails", product)
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

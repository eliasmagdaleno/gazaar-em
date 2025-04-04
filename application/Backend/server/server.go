package server

import (
	database "application/Database"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"application/Backend/routes"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func loadTemplate(filePath string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Printf("Error resolving absolute path for %s: %v", filePath, err)
		return "", err
	}
	log.Printf("Loading template from: %s", absPath)

	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Printf("Error reading file %s: %v", absPath, err)
		return "", err
	}

	return string(data), nil
}

func StartServer() {
	router := gin.Default()

	// Serve static files (e.g., images, CSS, JS)
	router.Static("/frontend", "../Frontend")

	// Trusted proxy configuration
	router.SetTrustedProxies([]string{"192.168.0.0/24"})

	routes.RegisterHomeRoutes(router)
	routes.RegisterMemberRoutes(router)
	routes.RegisterVPRoutes(router)
	routes.RegisterSearchRoutes(router)

	log.Println("🚀 Server running on http://0.0.0.0:8081")

	navbarPartial, err := loadTemplate("../Frontend/src/views/partials/navbar.hbs")
	if err != nil {
		log.Printf("Warning: Could not load navbar partial: %v", err)
	} else {
		raymond.RegisterPartial("navbar", navbarPartial)
	}

	eventCardPartial, err := loadTemplate("../Frontend/src/views/partials/eventcard.hbs")
	if err != nil {
		log.Printf("Warning: Could not load eventcard partial: %v", err)
	} else {
		raymond.RegisterPartial("eventcard", eventCardPartial)
	}

	productCardPartial, err := loadTemplate("../Frontend/src/views/partials/productcard.hbs")
	if err != nil {
		log.Printf("Warning: Could not load productcard partial: %v", err)
	} else {
		raymond.RegisterPartial("productcard", productCardPartial)
	}

	filterPartial, err := loadTemplate("../Frontend/src/views/partials/filter.hbs")
	if err != nil {
		log.Printf("Warning: Could not load filter partial: %v", err)
	} else {
		raymond.RegisterPartial("filter", filterPartial)
	}

	headerPartial, err := loadTemplate("../Frontend/src/views/partials/header.hbs")
	if err != nil {
		log.Printf("Warning: Could not load header partial: %v", err)
	} else {
		raymond.RegisterPartial("header", headerPartial)
	}

	// Register partials
	raymond.RegisterPartial("filter", filterPartial)
	raymond.RegisterPartial("header", headerPartial)

	// Example route to fetch search results from the database
	router.GET("/searchresults", func(c *gin.Context) {
		// Get the search string from the query parameters
		searchString := c.Query("q")
		if searchString == "" {
			c.String(http.StatusBadRequest, "Search query is required")
			return
		}

		// Execute the stored procedure
		rows, err := database.DB.Query("CALL search_items(?)", searchString)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error executing search query: %v", err))
			return
		}
		defer rows.Close()

		// Parse the results into events and products
		var events []map[string]string
		var products []map[string]string
		for rows.Next() {
			var itemType, thumbnail, title, host, conditionOrDate string
			if err := rows.Scan(&itemType, &thumbnail, &title, &host, &conditionOrDate); err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("Error scanning row: %v", err))
				return
			}

			if itemType == "event" {
				events = append(events, map[string]string{
					"thumbnail": thumbnail,
					"date":      conditionOrDate,
					"title":     title,
					"host":      host,
				})
			} else if itemType == "product" {
				products = append(products, map[string]string{
					"thumbnail": thumbnail,
					"title":     title,
					"host":      host,
					"condition": conditionOrDate,
				})
			}
		}

		// Load the layout template
		layoutTemplate, err := loadTemplate("../Frontend/src/views/layouts/layout.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout: %v", err))
			return
		}

		// Load the searchresults template
		searchResultsTemplate, err := loadTemplate("../Frontend/src/views/partials/searchresults.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading searchresults template: %v", err))
			return
		}

		// Render the searchresults content with data
		content, err := raymond.Render(searchResultsTemplate, map[string]interface{}{
			"title":    "Search Results",
			"events":   events,
			"products": products,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering searchresults content: %v", err))
			return
		}

		// Render the final layout with the content
		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "Search Results",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})

	log.Println("🚀 Server running on http://0.0.0.0:9081")
	router.Run("0.0.0.0:9081")
}

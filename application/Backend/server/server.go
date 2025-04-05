package server

import (
	//"database/sql"
	//"fmt"
	"log"
	//"net/http"
	"os"
	"path/filepath"

	"application/Backend/routes"
	//database "application/Database"

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

	// router.GET("/searchresults", func(c *gin.Context) {
	// 	// Get query parameters
	// 	category := c.DefaultQuery("category", "All")
	// 	searchValue := c.DefaultQuery("searchValue", "")

	// 	// Build the SQL query
	// 	var rows *sql.Rows
	// 	var err error
	// 	if category == "All" && searchValue == "" {
	// 		// No filters, show all items
	// 		rows, err = database.DB.Query("SELECT * FROM items")
	// 	} else if category == "All" {
	// 		// Only free text search
	// 		rows, err = database.DB.Query("SELECT * FROM items WHERE CONCAT(title, ' ', description) LIKE ?", "%"+searchValue+"%")
	// 	} else if searchValue == "" {
	// 		// Only category filter
	// 		rows, err = database.DB.Query("SELECT * FROM items WHERE category = ?", category)
	// 	} else {
	// 		// Both category and free text search
	// 		rows, err = database.DB.Query("SELECT * FROM items WHERE category = ? AND CONCAT(title, ' ', description) LIKE ?", category, "%"+searchValue+"%")
	// 	}

	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error executing search query: %v", err))
	// 		return
	// 	}
	// 	defer rows.Close()

	// 	// Parse the results
	// 	var items []map[string]string
	// 	for rows.Next() {
	// 		var id, title, description, category, thumbnail string
	// 		if err := rows.Scan(&id, &title, &description, &category, &thumbnail); err != nil {
	// 			c.String(http.StatusInternalServerError, fmt.Sprintf("Error scanning row: %v", err))
	// 			return
	// 		}

	// 		items = append(items, map[string]string{
	// 			"id":          id,
	// 			"title":       title,
	// 			"description": description,
	// 			"category":    category,
	// 			"thumbnail":   thumbnail,
	// 		})
	// 	}

	// 	// Load the layout and searchresults templates
	// 	layoutTemplate, err := loadTemplate("../Frontend/src/views/layouts/layout.hbs")
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout: %v", err))
	// 		return
	// 	}

	// 	searchResultsTemplate, err := loadTemplate("../Frontend/src/views/partials/searchresults.hbs")
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading searchresults template: %v", err))
	// 		return
	// 	}

	// 	// Render the searchresults content
	// 	content, err := raymond.Render(searchResultsTemplate, map[string]interface{}{
	// 		"items":       items,
	// 		"category":    category,
	// 		"searchValue": searchValue,
	// 		"totalCount":  len(items),
	// 	})
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering searchresults content: %v", err))
	// 		return
	// 	}

	// 	// Render the final layout
	// 	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
	// 		"title":   "Search Results",
	// 		"content": raymond.SafeString(content),
	// 	})
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
	// 		return
	// 	}

	// 	c.Header("Content-Type", "text/html")
	// 	c.String(http.StatusOK, output)
	// })

	log.Println("🚀 Server running on http://0.0.0.0:9081")
	router.Run("0.0.0.0:9081")
}

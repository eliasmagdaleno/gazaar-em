package routes

import (
	"application/Backend/core"
	"application/Backend/database"
	"application/Backend/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterMarketRoutes(router *gin.Engine) {
	// Apply the ProductCardMiddleware to the /market route
	router.GET("/market", RandomProductMiddleware(), UserIDMiddleware(),SignedInMiddleware(), marketHandler)
}

func fetchRandomMarketProducts(limit int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT image_url, title, price FROM items WHERE approve = 1 ORDER BY RAND() LIMIT %d", limit)
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var title, imageURL string
		var price float64
		err := rows.Scan(&imageURL, &title, &price)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		thumbnailPath := "frontend/assets/thumbnails/" + imageURL
		thumbnailFullPath := filepath.Join("assets/thumbnails", imageURL)
		originalImagePath := filepath.Join("assets/", imageURL)

		// Check if the thumbnail exists, if not, generate it
		if _, err := os.Stat(thumbnailFullPath); os.IsNotExist(err) {
			err := utils.GenerateThumbnail(originalImagePath, thumbnailFullPath, 150, 150) // Example size: 150x150
			if err != nil {
				log.Println("Thumbnail generation error:", err)
			}
		}

		products = append(products, map[string]interface{}{
			"title":     title,
			"price":     price,
			"thumbnail": thumbnailPath,
		})
	}
	return products, nil
}

func marketHandler(c *gin.Context) {
	// Retrieve the product cards set by the RandomProductMiddleware
	is_signed_in,_ := c.Get("is_signed_in")
	productCards, exists := c.Get("productCards")
	if !exists {
		c.String(http.StatusInternalServerError, "Error: Product cards not found in context")
		return
	}

	marketTemplate, err := core.LoadFrontendFile("src/views/market.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading market template: %v", err))
		return
	}

	content, err := raymond.Render(marketTemplate, map[string]interface{}{
		"title":    "Market",
		"products": productCards,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering market content: %v", err))
		return
	}

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Market",
		"content": raymond.SafeString(content),
		"is_signed_in": is_signed_in,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

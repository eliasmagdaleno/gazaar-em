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

func RegisterSearchRoutes(router *gin.Engine) {
	router.GET("/search", UserIDMiddleware(), SignedInMiddleware(), func(c *gin.Context) {
		q := c.Query("q")
		category := c.Query("category")
		is_signed_in, _ := c.Get("is_signed_in")
		userID, _ := c.Get("user_id")
		log.Println("User ID from context:", userID)


		query := "SELECT item_id, category, title, description, price, image_url FROM items WHERE 1=1 AND LOWER(category) != 'events' AND approve = 1"
		var args []interface{}

		if category != "" && category != "All" {
			query += " AND category = ?"
			args = append(args, category)
		}

		if q != "" {
			query += " AND CONCAT(title, ' ', description) LIKE ?"
			args = append(args, "%"+q+"%")
		}

		rows, err := database.DB.Query(query, args...)
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error executing search query: %v", err))
			return
		}
		defer rows.Close()

		var products []map[string]interface{}
		for rows.Next() {
			var itemID, category, title, description, imageURL string
			var price float64
			err := rows.Scan(&itemID, &category, &title, &description, &price, &imageURL)
			if err != nil {
				log.Println("Row scan error:", err)
				continue
			}

			thumbnailPath := "frontend/assets/thumbnails/" + imageURL
			thumbnailFullPath := filepath.Join("assets/thumbnails", imageURL)
			originalImagePath := filepath.Join("assets/", imageURL)

			// Check if the thumbnail exists, if not, generate it
			if _, err := os.Stat(thumbnailFullPath); os.IsNotExist(err) {
				err := utils.GenerateThumbnail(originalImagePath, thumbnailFullPath, 150, 150)
				if err != nil {
					log.Println("Thumbnail generation error:", err)
				}
			}

			products = append(products, map[string]interface{}{
				"item_id":          itemID,
				"category":    category,
				"title":       title,
				"description": description,
				"price":       price,
				"thumbnail":   thumbnailPath,
			})
		}

		data := map[string]interface{}{
			"category":      category,
			"q":             q,
			"count":         len(products),
			"products":      products,
			"isAll":         category == "All" || category == "",
			"isBooks":       category == "Books",
			"isElectronics": category == "Electronics",
			"isFurniture":   category == "Furniture",
		}

		searchResultsContent, err := core.LoadFrontendFile("src/views/partials/searchresults.hbs")
		if err != nil {
			log.Println("Template load error:", err)
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Template load error: %v", err))
			return
		}

		renderedSearchResults, err := raymond.Render(searchResultsContent, data)
		if err != nil {
			log.Println("Template render error:", err)
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Template render error: %v", err))
			return
		}

		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			log.Println("Layout load error:", err)
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Layout load error: %v", err))
			return
		}

		renderedLayout, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":         "Search Results",
			"content":       raymond.SafeString(renderedSearchResults),
			"q":             q,
			"category":      category,
			"isAll":         category == "All" || category == "",
			"isBooks":       category == "Books",
			"isElectronics": category == "Electronics",
			"isFurniture":   category == "Furniture",	
			"is_signed_in":  is_signed_in,
		})
		if err != nil {
			log.Println("Layout render error:", err)
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Layout render error: %v", err))
			return
		}

		// Send the rendered HTML as the response
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, renderedLayout)
	})
}

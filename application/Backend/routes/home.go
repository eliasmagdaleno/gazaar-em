package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"application/Backend/utils"
	"application/Backend/database"
	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

var (
	layoutTemplate string
	indexTemplate  string
)

func loadTemplate(filePath string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func RegisterHomeRoutes(r *gin.Engine) error {
	var err error

	// Serve static files from the Frontend/Assets folder
	r.Static("/assets", "Frontend/Assets")

	// Ensure thumbnails are generated
	thumbnails := []struct {
		inputPath  string
		outputPath string
	}{
		{"Frontend/Assets/event1.jpg", "Frontend/Assets/thumbnails/event1.jpg"},
		{"Frontend/Assets/event2.jpg", "Frontend/Assets/thumbnails/event2.jpg"},
		{"Frontend/Assets/event3.jpg", "Frontend/Assets/thumbnails/event3.jpg"},
		{"Frontend/Assets/event4.jpg", "Frontend/Assets/thumbnails/event4.jpg"},
		{"Frontend/Assets/product1.jpg", "Frontend/Assets/thumbnails/product1.jpg"},
		{"Frontend/Assets/product2.jpg", "Frontend/Assets/thumbnails/product2.jpg"},
		{"Frontend/Assets/product3.jpg", "Frontend/Assets/thumbnails/product3.jpg"},
		{"Frontend/Assets/product4.jpg", "Frontend/Assets/thumbnails/product4.jpg"},
	}

	for _, t := range thumbnails {
		err = utils.GenerateThumbnail(t.inputPath, t.outputPath, 150, 150)
		if err != nil {
			log.Printf("Error generating thumbnail for %s: %v", t.inputPath, err)
		}
	}

	// Load templates
	layoutTemplate, err = loadTemplate("Frontend/src/views/layouts/layout.hbs")
	if err != nil {
		return fmt.Errorf("error loading layout: %w", err)
	}

	indexTemplate, err = loadTemplate("Frontend/src/views/index.hbs")
	if err != nil {
		return fmt.Errorf("error loading index: %w", err)
	}

	r.GET("/", homeHandler)

	return nil
}

func fetchRandomProducts(limit int) ([]map[string]string,error){
	query := fmt.Sprintf("SELECT image_url, title, price FROM items ORDER BY RAND() LIMIT %d", limit)
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []map[string]string
	for rows.Next() {
		var thumbnail, title, price string
		err := rows.Scan(&thumbnail, &title, &price)
		if err != nil {
			return nil, err
		}
		products = append(products, map[string]string{
			"thumbnail": thumbnail,
			"title":     title,
			"price":     price,
		})
	}
	return products, nil
}

func homeHandler(c *gin.Context) {

	// events, err := fetchRandomEvents(4)
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching events: %v", err))
	// 	return
	// }

	products, err := fetchRandomProducts(4)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching products: %v", err))
		return
	}

	content, err := raymond.Render(indexTemplate, map[string]interface{}{
		"title":    "Home",
		"events":  []map[string]string{
			{"thumbnail": "/assets/thumbnails/event1.jpg", "date": "April 1, 2025", "title": "Event 1", "host": "Host A"},
			{"thumbnail": "/assets/thumbnails/event2.jpg", "date": "April 2, 2025", "title": "Event 2", "host": "Host B"},
			{"thumbnail": "/assets/thumbnails/event3.jpg", "date": "April 2, 2025", "title": "Event 3", "host": "Host C"},
			{"thumbnail": "/assets/thumbnails/event4.jpg", "date": "April 2, 2025", "title": "Event 4", "host": "Host D"},
		},
		"products": products,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering index content: %v", err))
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Home",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

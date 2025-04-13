package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
	"github.com/disintegration/imaging"
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

func generateThumbnail(inputPath, outputPath string, width, height int) error {
	// Open the source image
	src, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	// Resize the image to the specified dimensions
	thumbnail := imaging.Resize(src, width, height, imaging.Lanczos)

	// Save the thumbnail to the output path
	err = imaging.Save(thumbnail, outputPath)
	if err != nil {
		return fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return nil
}

func RegisterHomeRoutes(r *gin.Engine) error {
	var err error

	// Serve static files from the Frontend/Assets folder
	r.Static("/assets", "../Frontend/Assets")

	// Ensure thumbnails are generated
	thumbnails := []struct {
		inputPath  string
		outputPath string
	}{
		{"../Frontend/Assets/event1.jpg", "../Frontend/Assets/thumbnails/event1.jpg"},
		{"../Frontend/Assets/event2.jpg", "../Frontend/Assets/thumbnails/event2.jpg"},
		{"../Frontend/Assets/event3.jpg", "../Frontend/Assets/thumbnails/event3.jpg"},
		{"../Frontend/Assets/event4.jpg", "../Frontend/Assets/thumbnails/event4.jpg"},
		{"../Frontend/Assets/product1.jpg", "../Frontend/Assets/thumbnails/product1.jpg"},
		{"../Frontend/Assets/product2.jpg", "../Frontend/Assets/thumbnails/product2.jpg"},
		{"../Frontend/Assets/product3.jpg", "../Frontend/Assets/thumbnails/product3.jpg"},
		{"../Frontend/Assets/product4.jpg", "../Frontend/Assets/thumbnails/product4.jpg"},
	}

	for _, t := range thumbnails {
		err = generateThumbnail(t.inputPath, t.outputPath, 150, 150)
		if err != nil {
			log.Printf("Error generating thumbnail for %s: %v", t.inputPath, err)
		}
	}

	// Load templates
	layoutTemplate, err = loadTemplate("../Frontend/src/views/layouts/layout.hbs")
	if err != nil {
		return fmt.Errorf("error loading layout: %w", err)
	}

	indexTemplate, err = loadTemplate("../Frontend/src/views/index.hbs")
	if err != nil {
		return fmt.Errorf("error loading index: %w", err)
	}

	r.GET("/", homeHandler)

	return nil
}

func homeHandler(c *gin.Context) {
	content, err := raymond.Render(indexTemplate, map[string]interface{}{
		"title": "Home",
		"events": []map[string]string{
			{"thumbnail": "/assets/thumbnails/event1.jpg", "date": "April 1, 2025", "title": "Event 1", "host": "Host A"},
			{"thumbnail": "/assets/thumbnails/event2.jpg", "date": "April 2, 2025", "title": "Event 2", "host": "Host B"},
			{"thumbnail": "/assets/thumbnails/event3.jpg", "date": "April 2, 2025", "title": "Event 3", "host": "Host C"},
			{"thumbnail": "/assets/thumbnails/event4.jpg", "date": "April 2, 2025", "title": "Event 4", "host": "Host D"},
		},
		"products": []map[string]string{
			{"thumbnail": "/assets/thumbnails/product1.jpg", "title": "Product 1", "host": "$10", "condition": "New"},
			{"thumbnail": "/assets/thumbnails/product2.jpg", "title": "Product 2", "host": "$20", "condition": "Used"},
			{"thumbnail": "/assets/thumbnails/product3.jpg", "title": "Product 3", "host": "$15", "condition": "Used"},
			{"thumbnail": "/assets/thumbnails/product4.jpg", "title": "Product 4", "host": "$25", "condition": "Used"},
		},
		"searchResults": []map[string]string{
			{"thumbnail": "/assets/thumbnails/search1.jpg", "title": "Search Result 1", "description": "Description 1"},
			{"thumbnail": "/assets/thumbnails/search2.jpg", "title": "Search Result 2", "description": "Description 2"},
			{"thumbnail": "/assets/thumbnails/search3.jpg", "title": "Search Result 3", "description": "Description 3"},
			{"thumbnail": "/assets/thumbnails/search4.jpg", "title": "Search Result 4", "description": "Description 4"},
		},
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

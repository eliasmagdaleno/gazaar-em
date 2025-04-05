package routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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
			{"thumbnail": "/frontend/Assets/event1.jpg", "date": "April 1, 2025", "title": "Event 1", "host": "Host A"},
			{"thumbnail": "/frontend/Assets/event2.jpg", "date": "April 2, 2025", "title": "Event 2", "host": "Host B"},
			{"thumbnail": "/frontend/Assets/event2.jpg", "date": "April 2, 2025", "title": "Event 3", "host": "Host C"},
			{"thumbnail": "/frontend/Assets/event2.jpg", "date": "April 2, 2025", "title": "Event 4", "host": "Host D"},
		},
		"products": []map[string]string{
			{"thumbnail": "/frontend/Assets/product1.jpg", "title": "Product 1", "host": "$10", "condition": "New"},
			{"thumbnail": "/frontend/Assets/product2.jpg", "title": "Product 2", "host": "$20", "condition": "Used"},
			{"thumbnail": "/frontend/Assets/product2.jpg", "title": "Product 3", "host": "$15", "condition": "Used"},
			{"thumbnail": "/frontend/Assets/product2.jpg", "title": "Product 4", "host": "$25", "condition": "Used"},
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
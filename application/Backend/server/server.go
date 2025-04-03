package server

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

var frontendDir = os.Getenv("FRONTEND_PATH");


func loadFrontendFile(filelocation string) (string, error) {

	var sb strings.Builder
	sb.WriteString("../Frontend/")
	sb.WriteString(filelocation)

	abspath, err := filepath.Abs(sb.String())

	if err != nil {
		return "", err
	}


	data, err := os.ReadFile(abspath)

	// Error handling if path does not exsist
	if err != nil {
		return "", err
	}
	os.Stdout.Write(data)

	return string(data), nil
}

func TestDependency() {
	// Use raymond to prevent Go from removing it
	template := "Hello, {{name}}!"
	result, _ := raymond.Render(template, map[string]string{"name": "Zachary"})
	fmt.Println(result)
}

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

func StartServer() {
	router := gin.Default()

	router.Static("/frontend", "../Frontend")

	router.SetTrustedProxies([]string{"192.168.0.0/24"})

	// router.GET("/", func(c *gin.Context) {
	// 	htmlData, err := loadFrontendFile("src/html/index.html")
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
	// 		return
	// 	}

	// 	c.Header("Content-Type", "text/html")
	// 	c.String(http.StatusOK, string(htmlData))
	// })

	router.GET("/elias", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-elias.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/zachary", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-zachary.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/jiarong", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-jiarong.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/hemasri", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-hemasri.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/weiping", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-weiping.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	navbarPartial, err := loadTemplate("../Frontend/src/views/partials/navbar.hbs")
    if err != nil {
        log.Fatalf("Error loading navbar partial: %v", err)
    }
    eventCardPartial, err := loadTemplate("../Frontend/src/views/partials/eventcard.hbs")
    if err != nil {
        log.Fatalf("Error loading eventcard partial: %v", err)
    }
    productCardPartial, err := loadTemplate("../Frontend/src/views/partials/productcard.hbs")
    if err != nil {
        log.Fatalf("Error loading productcard partial: %v", err)
    }
	filterPartial, err := loadTemplate("../Frontend/src/views/partials/filter.hbs")
    if err != nil {
        log.Fatalf("Error loading filter partial: %v", err)
    }
	headerPartial, err := loadTemplate("../Frontend/src/views/partials/header.hbs")
    if err != nil {
        log.Fatalf("Error loading header partial: %v", err)
    }
	



        
        // Register partials
	raymond.RegisterPartial("navbar", navbarPartial)
    raymond.RegisterPartial("eventcard", eventCardPartial)
    raymond.RegisterPartial("productcard", productCardPartial)
	raymond.RegisterPartial("filter", filterPartial)
	raymond.RegisterPartial("header", headerPartial)


	router.GET("/", func(c *gin.Context) {
		layoutTemplate, err := loadTemplate("../Frontend/src/views/layouts/layout.hbs")
        if err != nil {
            c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout: %v", err))
            return
        }
		indexTemplate, err := loadTemplate("../Frontend/src/views/index.hbs")
        if err != nil {
            c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading template: %v", err))
            return
        }

	

        // Render the template with data
		content, err := raymond.Render(indexTemplate, map[string]interface{}{
            "title": "Home",
            "events": []map[string]string{
                {"thumbnail": "/frontend/Assets/event1.jpg", "date": "April 1, 2025", "title": "Event 1", "host": "Host A"},
                {"thumbnail": "/frontend/Assets/event2.jpg", "date": "April 2, 2025", "title": "Event 2", "host": "Host B"},
            },
            "products": []map[string]string{
                {"thumbnail": "/frontend/Assets/product1.jpg", "title": "Product 1", "host": "$10", "condition": "New"},
                {"thumbnail": "/frontend/Assets/product2.jpg", "title": "Product 2", "host": "$20", "condition": "Used"},
            },
        })
        if err != nil {
            c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering index content: %v", err))
            return
        }


		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
            "title": "Home",
            "content": raymond.SafeString(content),
        })
        if err != nil {
            c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
            return
        }

        c.Header("Content-Type", "text/html")
        c.String(http.StatusOK, output)
	})


	router.GET("/searchresults", func(c *gin.Context) {
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
			"title": "Search Results",
			"events": []map[string]string{
				{"thumbnail": "/frontend/Assets/event1.jpg", "date": "April 1, 2025", "title": "Event 1", "host": "Host A"},
				{"thumbnail": "/frontend/Assets/event2.jpg", "date": "April 2, 2025", "title": "Event 2", "host": "Host B"},
			},
			"products": []map[string]string{
				{"thumbnail": "/frontend/Assets/product1.jpg", "title": "Product 1", "host": "$10", "condition": "New"},
				{"thumbnail": "/frontend/Assets/product2.jpg", "title": "Product 2", "host": "$20", "condition": "Used"},
			},
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering searchresults content: %v", err))
			return
		}
	
		// Render the final layout with the content
		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title": "Search Results",
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

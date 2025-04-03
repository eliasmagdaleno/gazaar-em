package server

import (
	"log"

	"backend/routes"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	/*

		// Initialize the database connection
		if err := core.InitDB(); err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
	*/


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

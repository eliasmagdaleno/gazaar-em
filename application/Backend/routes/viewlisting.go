package routes

import (
	"fmt"
	"log"
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterViewListingsRoutes(router *gin.Engine) {
	router.GET("/viewlisting/:id", ProductDetailsMiddleware(), func(c *gin.Context) {
		log.Println("viewlisting: Entering viewlisting route")
		productDetails, exists := c.Get("productDetails")
		if !exists {
			log.Println("viewlisting: Product details not found in context")
			renderErrorPage(c, http.StatusInternalServerError, "Failed to load product details")
			return
		}

		log.Printf("viewlisting: Product details: %+v", productDetails) // Debugging log

		// Load the viewlisting.hbs template
		viewlistingTemplate, err := core.LoadFrontendFile("src/views/viewlisting.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading viewlisting template: %v", err))
			return
		}

		// Render the viewlisting.hbs template with product details
		content, err := raymond.Render(viewlistingTemplate, productDetails)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering viewlisting template: %v", err))
			return
		}

		// Load the layout.hbs template
		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			return
		}

		// Render the layout.hbs template with the viewlisting content
		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "View Listing",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})
}

func RegisterCreateListingRoutes(router *gin.Engine) {
	router.GET("/createlisting", createListingHandler)
	router.POST("/createlisting", submitListingHandler)

	// New route for selectlocation
	router.GET("/selectlocation", selectLocationHandler)
}

func createListingHandler(c *gin.Context) {
	createListingTemplate, err := core.LoadFrontendFile("src/views/createlisting.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading create listing template: %v", err))
		return
	}

	content, err := raymond.Render(createListingTemplate, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering create listing template: %v", err))
		return
	}

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Create Listing",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

func submitListingHandler(c *gin.Context) {
	// Handle form submission logic here
	c.String(http.StatusOK, "Listing submitted successfully!")
}

// Handler for selectlocation page
func selectLocationHandler(c *gin.Context) {
	selectLocationTemplate, err := core.LoadFrontendFile("src/views/selectlocation.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading select location template: %v", err)
		return
	}

	content, err := raymond.Render(selectLocationTemplate, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering select location template: %v", err)
		return
	}

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading layout template: %v", err)
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Select Location",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering layout: %v", err)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

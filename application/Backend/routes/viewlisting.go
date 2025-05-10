package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterViewListingsRoutes(router *gin.Engine) {
	router.GET("/listing", func(c *gin.Context) {
		viewListingTemplate, err := core.LoadFrontendFile("src/views/viewlisting.hbs")
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}
		content, err := raymond.Render(viewListingTemplate, map[string]interface{}{
			"title":     "View Listings",
			"thumbnail": "assets/thumbnails/event1.jpg",
			"fullImage": "assets/thumbnails/event1.jpg",
		})
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error rendering template: %v", err))
			return
		}
		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			return
		}
		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "View Listing",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})
}

func RegisterCreateListingRoutes(router *gin.Engine) {
	router.GET("/createlisting", createListingHandler)
	router.POST("/createlisting", submitListingHandler)
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

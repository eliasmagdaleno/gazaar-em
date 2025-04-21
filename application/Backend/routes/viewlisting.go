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
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}
		content, err := raymond.Render(viewListingTemplate, map[string]interface{}{
			"title": "View Listings",
			"thumbnail": "assets/thumbnails/event1.jpg",
			"fullImage": "assets/thumbnails/event1.jpg",
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering template: %v", err))
			return
		}
		 // Load the layout template
		 layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		 if err != nil {
			 c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			 return
		 }
 
		 // Render the layout with the viewlisting content
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

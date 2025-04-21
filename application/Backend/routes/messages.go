package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"
	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterMessagesRoutes(router *gin.Engine) {
	router.GET("/messages", func(c *gin.Context) {
		messagesTemplate, err := core.LoadFrontendFile("src/views/messages.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}
		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		 if err != nil {
			 c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			 return
		 }
		 content, err := raymond.Render(messagesTemplate, map[string]interface{}{
			"title": "View Messages",
			
		})
 
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
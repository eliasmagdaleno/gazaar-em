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
		// Test data for messages
		messages := []map[string]interface{}{
			{"isSender": true, "message": "Hi, is the item still available?", "timestamp": "10:00 AM"},
			{"isSender": false, "message": "Yes, it is available.", "timestamp": "10:02 AM"},
			{"isSender": true, "message": "Great! Can I pick it up tomorrow?", "timestamp": "10:05 AM"},
			{"isSender": false, "message": "Sure, what time works for you?", "timestamp": "10:07 AM"},
		}

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
			"title":    "View Messages",
			"messages": messages,
		})

		// Render the layout with the messages content
		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "View Messages",
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

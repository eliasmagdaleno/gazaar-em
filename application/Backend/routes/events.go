package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterEventsRoutes(router *gin.Engine) {
	router.GET("/events", eventsHandler)
}

func eventsHandler(c *gin.Context) {
	// Fake events data
	events := []map[string]string{
		{"thumbnail": "/assets/thumbnails/event1.jpg", "date": "May 1, 2025", "title": "Event 1", "host": "Host A", "location": "San Francisco"},
		{"thumbnail": "/assets/thumbnails/event2.jpg", "date": "May 2, 2025", "title": "Event 2", "host": "Host B", "location": "Los Angeles"},
		{"thumbnail": "/assets/thumbnails/event3.jpg", "date": "May 3, 2025", "title": "Event 3", "host": "Host C", "location": "New York"},
		{"thumbnail": "/assets/thumbnails/event4.jpg", "date": "May 4, 2025", "title": "Event 4", "host": "Host D", "location": "Chicago"},
	}

	eventsTemplate, err := core.LoadFrontendFile("src/views/events.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading events template: %v", err))
		return
	}

	content, err := raymond.Render(eventsTemplate, map[string]interface{}{
		"title":  "Events",
		"events": events,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering events content: %v", err))
		return
	}

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Events",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

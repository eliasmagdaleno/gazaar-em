package routes

import (
	"fmt"
	"log"
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterEventsRoutes(router *gin.Engine) {
	// Apply the EventCardMiddleware to the /events route
	router.GET("/events", RandomEventMiddleware(), eventsHandler)
}

func eventsHandler(c *gin.Context) {
	// Temporarily replace the eventCards data with a simplified hardcoded structure for testing
	hardcodedEvents := []map[string]interface{}{
		{"thumbnail": "frontend/assets/thumbnails/test1.jpg", "title": "Test Event 1", "postDate": "2025-05-01"},
		{"thumbnail": "frontend/assets/thumbnails/test2.jpg", "title": "Test Event 2", "postDate": "2025-05-02"},
	}
	log.Printf("[DEBUG] Using hardcoded events for testing: %+v", hardcodedEvents)

	eventsTemplate, err := core.LoadFrontendFile("src/views/events.hbs")
	if err != nil {
		log.Printf("eventsHandler: Error loading events template: %v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading events template: %v", err))
		return
	}

	// Add detailed logging to trace template rendering
	log.Printf("[DEBUG] Rendering events.hbs with data: %+v", map[string]interface{}{
		"title":  "Events",
		"events": hardcodedEvents,
	})

	content, err := raymond.Render(eventsTemplate, map[string]interface{}{
		"title":  "Events",
		"events": hardcodedEvents,
	})
	if err != nil {
		log.Printf("[ERROR] Error rendering events.hbs: %v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering events.hbs: %v", err))
		return
	}

	// Add logging to confirm the events data is accessible in the template
	log.Printf("[DEBUG] Retrieved events data: %+v", hardcodedEvents)

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		log.Printf("eventsHandler: Error loading layout template: %v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
		return
	}

	// Add detailed logging to trace template rendering
	log.Printf("[DEBUG] Rendering layout.hbs with content: %s", content)

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Events",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		log.Printf("[ERROR] Error rendering layout.hbs: %v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout.hbs: %v", err))
		return
	}

	log.Printf("eventsHandler: Successfully rendered events page")
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

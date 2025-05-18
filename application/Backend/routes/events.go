package routes

import (
	"fmt"
	"log"
	"net/http"

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterEventsRoutes(router *gin.Engine) {
	// Apply the EventCardMiddleware to the /events route
	router.GET("/events", RandomEventMiddleware(), eventsHandler)
}

func eventsHandler(c *gin.Context) {

	rows, err := database.DB.Query(
		`SELECT item_id, image_url, title, seller_id,
		DATE_FORMAT(post_date, '%M %e, %Y') AS date,
		description
   FROM items
  	ORDER BY post_date DESC`,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("DB error: %v", err))
		return
	}
	defer rows.Close()

	// Build a slice of event maps
	var evs []map[string]interface{}
	for rows.Next() {
		var id int
		var img, title, host, date, desc string
		if err := rows.Scan(&id, &img, &title, &host, &date, &desc); err != nil {
			continue
		}
		evs = append(evs, map[string]interface{}{
			"id":        id,
			"thumbnail": "/assets/thumbnails/" + img,
			"title":     title,
			"host":      host,
			"date":      date,
			"location":  desc,
		})
	}

	eventsTemplate, err := core.LoadFrontendFile("src/views/events.hbs")
	if err != nil {
		log.Printf("eventsHandler: Error loading events template: %v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading events template: %v", err))
		return
	}

	// Add detailed logging to trace template rendering
	log.Printf("[DEBUG] Rendering events.hbs with data: %+v", map[string]interface{}{
		"title":  "Events",
		"events": evs,
	})

	content, err := raymond.Render(eventsTemplate, map[string]interface{}{
		"title":  "Events",
		"events": evs,
	})
	if err != nil {
		log.Printf("[ERROR] Error rendering events.hbs: %v", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering events.hbs: %v", err))
		return
	}

	// Add logging to confirm the events data is accessible in the template
	log.Printf("[DEBUG] Retrieved events data: %+v", evs)

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

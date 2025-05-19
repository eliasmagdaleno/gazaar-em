package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterEventsRoutes(router *gin.Engine) {
	// Apply the EventCardMiddleware to the /events route
	router.GET("/events", RandomEventMiddleware(), eventsHandler)
	router.POST("/events/delete", deleteEventHandler)
}

func eventsHandler(c *gin.Context) {

	rows, err := database.DB.Query(
	`SELECT items.item_id, items.image_url, items.title, Account.user_name AS host,
       DATE_FORMAT(items.post_date, '%M %e, %Y') AS date,
       items.description
		FROM items
		JOIN Account ON items.seller_id = Account.user_id
		WHERE LOWER(items.category) = 'events' AND items.approve = 1
		ORDER BY items.post_date DESC`,
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
			"thumbnail": "../../frontend/assets/thumbnails/" + img,
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
	// log.Printf("[DEBUG] Rendering events.hbs with data: %+v", map[string]interface{}{
	// 	"title":  "Events",
	// 	"events": evs,
	// })

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

func deleteEventHandler(c *gin.Context) {
	// 1) Parse the item ID
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}

	// 2) Look up the seller_id for this item
	var sellerID int
	err = database.DB.
		QueryRow("SELECT seller_id FROM items WHERE item_id = ?", id).
		Scan(&sellerID)
	if err == sql.ErrNoRows {
		c.String(http.StatusNotFound, "Item not found")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("DB error: %v", err))
		return
	}

	// 3) Delete in a transaction: first remove child rows
	tx, err := database.DB.Begin()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to begin transaction")
		return
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		"DELETE FROM transactions WHERE seller_id = ?", sellerID,
	); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed deleting transactions: %v", err))
		return
	}

	// 4) Now delete the item
	if _, err := tx.Exec(
		"DELETE FROM items WHERE item_id = ?", id,
	); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed deleting item: %v", err))
		return
	}

	if err := tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Commit error: %v", err))
		return
	}

	// 5) Redirect back
	c.Redirect(http.StatusSeeOther, "/events")
}
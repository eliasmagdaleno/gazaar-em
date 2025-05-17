package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterMessagesRoutes(router *gin.Engine) {
	// GET endpoint to view messages
	router.GET("/messages", func(c *gin.Context) {
		roomId := c.Query("room")
		var rows *sql.Rows
		var err error
		if roomId != "" {
			rows, err = database.DB.Query(`
				SELECT content,
					   DATE_FORMAT(timestamp, '%Y-%m-%d %H:%i') AS ts
				  FROM Message
				 WHERE room = ?
				ORDER BY timestamp ASC
			`, roomId)
		} else {
			rows, err = database.DB.Query(`
				SELECT content,
					   DATE_FORMAT(timestamp, '%Y-%m-%d %H:%i') AS ts
				  FROM Message
				ORDER BY timestamp ASC
			`)
		}
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("DB error: %v", err))
			return
		}
		defer rows.Close()

		// Build a flat slice of message maps
		var allMsgs []map[string]string
		for rows.Next() {
			var text, ts string
			if err := rows.Scan(&text, &ts); err != nil {
				continue
			}
			allMsgs = append(allMsgs, map[string]string{
				"message":   text,
				"timestamp": ts,
			})
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
			"messages": allMsgs,
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

	// POST endpoint to send messages
	router.POST("/messages", func(c *gin.Context) {
		message := c.PostForm("message")
		if message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message cannot be empty"})
			return
		}

		// Insert message into database (no room)
		_, err := database.DB.Exec(`
			INSERT INTO Message (content, timestamp)
			VALUES (?, NOW())
		`, message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save message: %v", err)})
			return
		}

		// Redirect to the messages page
		c.Redirect(http.StatusSeeOther, "/messages")
	})
}

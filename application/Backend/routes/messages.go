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

		// Fetch all unique room IDs (sellerIDs) from the Message table
		roomRows, err := database.DB.Query("SELECT DISTINCT room FROM Message")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("DB error: %v", err))
			return
		}
		defer roomRows.Close()

		var rooms []map[string]interface{}
		for roomRows.Next() {
			var roomId int
			if err := roomRows.Scan(&roomId); err != nil {
				continue
			}
			// Optionally, fetch the seller's name from a Users table if you have one
			var sellerName string
			err = database.DB.QueryRow("SELECT name FROM Users WHERE id = ?", roomId).Scan(&sellerName)
			if err != nil {
				sellerName = fmt.Sprintf("User %d", roomId)
			}
			// Optionally, fetch the last message for this room
			var lastMessage string
			_ = database.DB.QueryRow("SELECT content FROM Message WHERE room = ? ORDER BY timestamp DESC LIMIT 1", roomId).Scan(&lastMessage)

			rooms = append(rooms, map[string]interface{}{
				"id":          roomId,
				"name":        sellerName,
				"lastMessage": lastMessage,
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
			"rooms":    rooms,
			"roomId":   roomId,
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
		roomId := c.PostForm("room")
		if message == "" || roomId == "" {
			// Just redirect back to the current chat room (or /messages if no room)
			if roomId != "" {
				c.Redirect(http.StatusSeeOther, "/messages?room="+roomId)
			} else {
				c.Redirect(http.StatusSeeOther, "/messages")
			}
			return
		}

		// Insert message into database with room
		_, err := database.DB.Exec(`
			INSERT INTO Message (content, timestamp, room)
			VALUES (?, NOW(), ?)
		`, message, roomId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save message: %v", err)})
			return
		}

		// Redirect to the seller's chat room
		c.Redirect(http.StatusSeeOther, "/messages?room="+roomId)
	})

	// POST endpoint to delete a chat room
	router.POST("/messages/delete-room", func(c *gin.Context) {
		roomId := c.PostForm("room")
		if roomId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
			return
		}

		// Delete all messages in this room
		_, err := database.DB.Exec("DELETE FROM Message WHERE room = ?", roomId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete messages: %v", err)})
			return
		}

		// Optionally, redirect to /messages (no room selected)
		c.Redirect(http.StatusSeeOther, "/messages")
	})
}

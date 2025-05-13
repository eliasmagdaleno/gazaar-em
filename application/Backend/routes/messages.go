package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterMessagesRoutes(router *gin.Engine) {
	router.GET("/messages", func(c *gin.Context) {

		rows, err := database.DB.Query(`
      SELECT content,
             DATE_FORMAT(timestamp, '%Y-%m-%d %H:%i') AS ts
        FROM Message
       ORDER BY timestamp ASC
    `)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("DB error: %v", err))
			return
		}
		defer rows.Close()

		// 2) Build a flat slice of message maps
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
}

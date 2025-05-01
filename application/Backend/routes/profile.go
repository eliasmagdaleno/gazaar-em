package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterProfileRoutes(router *gin.Engine) {
	router.GET("/profile", profileHandler)
}

func profileHandler(c *gin.Context) {
	// Fetch user information (replace with actual user ID or session data)
	//userID := 1 // Example user ID
	username := "John Doe"
	bio := "Welcome to my profile! I love selling and hosting events."

	// Fetch user's posted items
	// itemsQuery := "SELECT thumbnail, title, price, condition FROM items WHERE user_id = ?"
	// itemsRows, err := database.DB.Query(itemsQuery, userID)
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching items: %v", err))
	// 	return
	// }
	// defer itemsRows.Close()

	var items []map[string]string
	// for itemsRows.Next() {
	// 	var thumbnail, title, condition string
	// 	var price float64
	// 	if err := itemsRows.Scan(&thumbnail, &title, &price, &condition); err != nil {
	// 		continue
	// 	}
	// 	items = append(items, map[string]string{
	// 		"thumbnail": thumbnail,
	// 		"title":     title,
	// 		"price":     fmt.Sprintf("%.2f", price),
	// 		"condition": condition,
	// 	})
	// }

	// Fetch user's posted events
	// eventsQuery := "SELECT thumbnail, date, title, host FROM events WHERE user_id = ?"
	// eventsRows, err := database.DB.Query(eventsQuery, userID)
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching events: %v", err))
	// 	return
	// }
	// defer eventsRows.Close()

	var events []map[string]string
	// for eventsRows.Next() {
	// 	var thumbnail, date, title, host string
	// 	if err := eventsRows.Scan(&thumbnail, &date, &title, &host); err != nil {
	// 		continue
	// 	}
	// 	events = append(events, map[string]string{
	// 		"thumbnail": thumbnail,
	// 		"date":      date,
	// 		"title":     title,
	// 		"host":      host,
	// 	})
	// }

	// Render the profile page
	profileTemplate, err := core.LoadFrontendFile("src/views/profile.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading profile template: %v", err))
		return
	}

	content, err := raymond.Render(profileTemplate, map[string]interface{}{
		"username": username,
		"bio":      bio,
		"items":    items,
		"events":   events,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering profile content: %v", err))
		return
	}

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "User Profile",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

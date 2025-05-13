package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterProfileRoutes(router *gin.Engine) {
	router.GET("/profile", func(c *gin.Context) {

		userID := c.GetInt("user_id")

		// 2) Fetch user info (username)
		var username string
		err := database.DB.QueryRow(
			`SELECT user_name
           FROM Account
          WHERE user_id = ?`,
			userID,
		).Scan(&username)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching user info: %v", err))
			return
		}

		// 3) Fetch items posted by user
		itemRows, err := database.DB.Query(
			`SELECT item_id, image_url, title, price,
                DATE_FORMAT(post_date, '%M %e, %Y') AS postDate
           FROM items
          WHERE seller_id = ?
          ORDER BY post_date DESC`,
			userID,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching items: %v", err))
			return
		}
		defer itemRows.Close()

		var items []map[string]interface{}
		for itemRows.Next() {
			var id int
			var img, title, postDate string
			var price float64
			if err := itemRows.Scan(&id, &img, &title, &price, &postDate); err != nil {
				continue
			}
			items = append(items, map[string]interface{}{
				"id":       id,
				"imageUrl": "/assets/thumbnails/" + img,
				"title":    title,
				"price":    price,
				"postDate": postDate,
			})
		}

		// 4) Fetch "events" (category == 'others')
		evtRows, err := database.DB.Query(
			`SELECT item_id, image_url, title,
                DATE_FORMAT(post_date, '%M %e, %Y') AS postDate,
                description
           FROM items
          WHERE seller_id = ?
          ORDER BY post_date DESC`,
			userID,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching events: %v", err))
			return
		}
		defer evtRows.Close()

		var events []map[string]interface{}
		for evtRows.Next() {
			var id int
			var img, title, postDate, desc string
			if err := evtRows.Scan(&id, &img, &title, &postDate, &desc); err != nil {
				continue
			}
			events = append(events, map[string]interface{}{
				"id":        id,
				"thumbnail": "/assets/thumbnails/" + img,
				"title":     title,
				"date":      postDate,
				"location":  desc,
			})
		}

		// 5) Fetch reviews where this user is seller
		revRows, err := database.DB.Query(
			`SELECT review_id, buyer_id, rating, review_text,
                DATE_FORMAT(review_date, '%M %e, %Y') AS reviewDate
           FROM Review
          WHERE seller_id = ?
          ORDER BY review_date DESC`,
			userID,
		)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching reviews: %v", err))
			return
		}
		defer revRows.Close()

		var reviews []map[string]interface{}
		for revRows.Next() {
			var rid, buyerID, rating int
			var text, reviewDate string
			if err := revRows.Scan(&rid, &buyerID, &rating, &text, &reviewDate); err != nil {
				continue
			}
			// lookup buyer name
			var buyerName string
			database.DB.QueryRow(`SELECT user_name FROM Account WHERE user_id = ?`, buyerID).Scan(&buyerName)

			reviews = append(reviews, map[string]interface{}{
				"reviewID":  rid,
				"buyerName": buyerName,
				"rating":    rating,
				"text":      text,
				"date":      reviewDate,
			})
		}

		profileTemplate, err := core.LoadFrontendFile("src/views/profile.hbs")
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error loading profile template: %v", err))
			return
		}

		content, err := raymond.Render(profileTemplate, map[string]interface{}{
			"username": username,
			"items":    items,
		})
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error rendering profile content: %v", err))
			return
		}

		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			return
		}

		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "Profile",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			renderErrorPage(c, http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})
}

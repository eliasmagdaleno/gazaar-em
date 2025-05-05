package routes

import (
	"fmt"
	"net/http"

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterMarketRoutes(router *gin.Engine) {
	router.GET("/market", marketHandler)
}

func fetchRandomMarketProducts(limit int) ([]map[string]string, error) {
	query := fmt.Sprintf("SELECT image_url, title, price FROM items ORDER BY RAND() LIMIT %d", limit)
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []map[string]string
	for rows.Next() {
		var thumbnail, title, price string
		err := rows.Scan(&thumbnail, &title, &price)
		if err != nil {
			return nil, err
		}
		products = append(products, map[string]string{
			"thumbnail": thumbnail,
			"title":     title,
			"price":     price,
		})
	}
	return products, nil
}

func marketHandler(c *gin.Context) {
	products, err := fetchRandomMarketProducts(20)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching products: %v", err))
		return
	}

	marketTemplate, err := core.LoadFrontendFile("src/views/market.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading market template: %v", err))
		return
	}

	content, err := raymond.Render(marketTemplate, map[string]interface{}{
		"title":    "Market",
		"products": products,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering market content: %v", err))
		return
	}

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Market",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

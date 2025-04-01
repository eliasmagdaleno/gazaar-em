package routes

import (
	"fmt"
	"log"
	"net/http"

	"backend/core"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

type Item struct {
	ID          int
	Category    string
	Title       string
	Description string
	Price       float64
	ImageFull   string
	ImageThumb  string
}

func RegisterSearchRoutes(router *gin.Engine) {
	router.GET("/search", searchHandler)
}

func searchHandler(c *gin.Context) {
	category := c.Query("category")
	q := c.Query("q")

	query := "SELECT id, category, title, description, price, image_full, image_thumb FROM items WHERE 1=1"
	var args []interface{}
	if category != "" && category != "all" {
		query += " AND category = ?"
		args = append(args, category)
	}
	if q != "" {
		query += " AND CONCAT(title, ' ', description) LIKE ?"
		args = append(args, "%"+q+"%")
	}

	rows, err := core.DB.Query(query, args...)
	if err != nil {
		log.Println("DB Query error:", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("DB Query error: %v", err))
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var it Item
		err := rows.Scan(&it.ID, &it.Category, &it.Title, &it.Description, &it.Price, &it.ImageFull, &it.ImageThumb)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		items = append(items, it)
	}

	data := map[string]interface{}{
		"category": category,
		"q":        q,
		"count":    len(items),
		"items":    items,
	}

	tmpl, err := core.LoadFrontendFile("src/html/search_results.hbs")
	if err != nil {
		log.Println("Template load error:", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Template load error: %v", err))
		return
	}

	rendered, err := raymond.Render(tmpl, data)
	if err != nil {
		log.Println("Template render error:", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Template render error: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, rendered)
}

package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

var frontendDir = os.Getenv("FRONTEND_PATH")
var db *sql.DB

type Item struct {
	ID          int
	Category    string
	Title       string
	Description string
	Price       float64
	ImageFull   string
	ImageThumb  string
}

func loadFrontendFile(filelocation string) (string, error) {

	var sb strings.Builder
	sb.WriteString("../Frontend/")
	sb.WriteString(filelocation)

	abspath, err := filepath.Abs(sb.String())

	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(abspath)

	// Error handling if path does not exsist
	if err != nil {
		return "", err
	}
	os.Stdout.Write(data)

	return string(data), nil
}

func TestDependency() {
	// Use raymond to prevent Go from removing it
	template := "Hello, {{name}}!"
	result, _ := raymond.Render(template, map[string]string{"name": "Zachary"})
	fmt.Println(result)
}

// searchHandler processes the search request from the VP test home page.
func searchHandler(c *gin.Context) {
	// Retrieve query parameters.
	category := c.Query("category")
	q := c.Query("q")

	// Build the SQL query.
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

	// Execute the query.
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println("DB Query error:", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("DB Query error: %v", err))
		return
	}
	defer rows.Close()

	// Collect the items.
	var items []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.Category, &it.Title, &it.Description, &it.Price, &it.ImageFull, &it.ImageThumb); err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		items = append(items, it)
	}

	// Prepare data for the template.
	data := map[string]interface{}{
		"category": category,
		"q":        q,
		"count":    len(items),
		"items":    items,
	}

	// Load the search results Handlebars template.
	tmpl, err := loadFrontendFile("src/html/search_results.hbs")
	if err != nil {
		log.Println("Template load error:", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Template load error: %v", err))
		return
	}

	// Render the template with data.
	rendered, err := raymond.Render(tmpl, data)
	if err != nil {
		log.Println("Template render error:", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Template render error: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, rendered)
}

func StartServer() {

	router := gin.Default()

	router.Static("/frontend", "../Frontend")

	router.SetTrustedProxies([]string{"192.168.0.0/24"})

	router.GET("/", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/elias", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-elias.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/zachary", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-zachary.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/jiarong", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-jiarong.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/hemasri", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-hemasri.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	router.GET("/weiping", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/member-pages/aboutme-weiping.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, string(htmlData))
	})

	// VP Test Home Page (search form).
	router.GET("/vp", func(c *gin.Context) {
		htmlData, err := loadFrontendFile("src/html/vp_home.html")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlData)
	})

	// New search route.
	router.GET("/search", searchHandler)

	log.Println("🚀 Server running on http://0.0.0.0:8081")
	router.Run("0.0.0.0:8081")
}

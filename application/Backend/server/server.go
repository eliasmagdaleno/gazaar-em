package server

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

var frontendDir = os.Getenv("FRONTEND_PATH");


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



	

	log.Println("🚀 Server running on http://0.0.0.0:8081")
	router.Run("0.0.0.0:8081")
}

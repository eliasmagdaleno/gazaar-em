package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func TestDependency() {
	// Use raymond to prevent Go from removing it
	template := "Hello, {{name}}!"
	result, _ := raymond.Render(template, map[string]string{"name": "Zachary"})
	fmt.Println(result)
}

func StartServer() {
	router := gin.Default()

	router.SetTrustedProxies([]string{"192.168.0.0/24"})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Go!")
	})

	log.Println("🚀 Server running on http://0.0.0.0:8081")
	router.Run("0.0.0.0:8081")
}

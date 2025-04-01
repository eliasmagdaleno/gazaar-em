package server

import (
	"log"

	"backend/core"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	// Initialize the database connection
	if err := core.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	router := gin.Default()

	// Serve static files (e.g., images, CSS, JS)
	router.Static("/frontend", "../Frontend")

	// Trusted proxy configuration
	router.SetTrustedProxies([]string{"192.168.0.0/24"})

	// Register all routes from routes package
	routes.RegisterHomeRoutes(router)
	routes.RegisterMemberRoutes(router)
	routes.RegisterVPRoutes(router)
	routes.RegisterSearchRoutes(router)

	log.Println("🚀 Server running on http://0.0.0.0:8081")
	router.Run("0.0.0.0:8081")
}

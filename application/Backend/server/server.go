package server

import (
	//"database/sql"
	//"fmt"
	"log"
	//"net/http"
	"os"
	"path/filepath"

	"application/Backend/routes"
	//database "application/Database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func loadTemplate(filePath string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Printf("Error resolving absolute path for %s: %v", filePath, err)
		return "", err
	}
	log.Printf("Loading template from: %s", absPath)

	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Printf("Error reading file %s: %v", absPath, err)
		return "", err
	}

	return string(data), nil
}

func StartServer() {
	router := gin.Default()

	// Serve static files (e.g., images, CSS, JS)
	router.Static("/frontend", "../Frontend")

	// Trusted proxy configuration
	router.SetTrustedProxies([]string{"192.168.0.0/24"})

	routes.RegisterHomeRoutes(router)
	routes.RegisterMemberRoutes(router)
	routes.RegisterVPRoutes(router)
	routes.RegisterSearchRoutes(router)

	log.Println("🚀 Server running on http://0.0.0.0:8081")

	navbarPartial, err := loadTemplate("../Frontend/src/views/partials/navbar.hbs")
	if err != nil {
		log.Printf("Warning: Could not load navbar partial: %v", err)
	} else {
		raymond.RegisterPartial("navbar", navbarPartial)
	}

	eventCardPartial, err := loadTemplate("../Frontend/src/views/partials/eventcard.hbs")
	if err != nil {
		log.Printf("Warning: Could not load eventcard partial: %v", err)
	} else {
		raymond.RegisterPartial("eventcard", eventCardPartial)
	}

	productCardPartial, err := loadTemplate("../Frontend/src/views/partials/productcard.hbs")
	if err != nil {
		log.Printf("Warning: Could not load productcard partial: %v", err)
	} else {
		raymond.RegisterPartial("productcard", productCardPartial)
	}

	filterPartial, err := loadTemplate("../Frontend/src/views/partials/filter.hbs")
	if err != nil {
		log.Printf("Warning: Could not load filter partial: %v", err)
	} else {
		raymond.RegisterPartial("filter", filterPartial)
	}

	headerPartial, err := loadTemplate("../Frontend/src/views/partials/header.hbs")
	if err != nil {
		log.Printf("Warning: Could not load header partial: %v", err)
	} else {
		raymond.RegisterPartial("header", headerPartial)
	}

	

	log.Println("🚀 Server running on http://0.0.0.0:9081")
	router.Run("0.0.0.0:9081")
}

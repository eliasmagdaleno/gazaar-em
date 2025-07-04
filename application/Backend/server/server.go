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
	router.Static("/frontend", "Frontend")
	router.Static("/assets", "Frontend/assets")

	// Trusted proxy configuration
	router.SetTrustedProxies([]string{"192.168.0.0/24"})
	
	router.Use(routes.UserIDMiddleware())
	router.Use(routes.SignedInMiddleware())
	



	routes.RegisterHomeRoutes(router)
	routes.RegisterMemberRoutes(router)
	routes.RegisterVPRoutes(router)
	routes.RegisterSearchRoutes(router)
	routes.RegisterViewListingsRoutes(router)
	routes.RegisterMessagesRoutes(router)
	routes.RegisterAuthRoutes(router)
	routes.RegisterProfileRoutes(router)
	routes.RegisterMarketRoutes(router)
	routes.RegisterEventsRoutes(router)
	routes.RegisterCreateListingRoutes(router)

	log.Println("🚀 Server running on http://0.0.0.0:8081")

	navbarPartial, err := loadTemplate("Frontend/src/views/partials/navbar.hbs")
	if err != nil {
		log.Printf("Warning: Could not load navbar partial: %v", err)
	} else {
		raymond.RegisterPartial("navbar", navbarPartial)
	}

	eventCardPartial, err := loadTemplate("Frontend/src/views/partials/eventcard.hbs")
	if err != nil {
		log.Printf("Warning: Could not load eventcard partial: %v", err)
	} else {
		raymond.RegisterPartial("eventcard", eventCardPartial)
	}

	productCardHomePartial, err := loadTemplate("Frontend/src/views/partials/productcard-home.hbs")
	if err != nil {
		log.Printf("Warning: Could not load productcard-home partial: %v", err)
	} else {
		raymond.RegisterPartial("productcard-home", productCardHomePartial)
	}

	productCardSrPartial, err := loadTemplate("Frontend/src/views/partials/productcard-search.hbs")
	if err != nil {
		log.Printf("Warning: Could not load productcard-search partial: %v", err)
	} else {
		raymond.RegisterPartial("productcard-search", productCardSrPartial)
	}

	productCardMarketplacePartial, err := loadTemplate("Frontend/src/views/partials/productcard-marketplace.hbs")
	if err != nil {
		log.Printf("Warning: Could not load productcard-marketplace partial: %v", err)
	} else {
		raymond.RegisterPartial("productcard-marketplace", productCardMarketplacePartial)
	}

	eventCardMarketplacePartial, err := loadTemplate("Frontend/src/views/partials/eventcard-marketplace.hbs")
	if err != nil {
		log.Printf("Warning: Could not load eventcard-marketplace partial: %v", err)
	} else {
		raymond.RegisterPartial("eventcard-marketplace", eventCardMarketplacePartial)
	}

	senderCardPartial, err := loadTemplate("Frontend/src/views/partials/sendercard.hbs")
	if err != nil {
		log.Printf("Warning: Could not load sendercard partial: %v", err)
	} else {
		raymond.RegisterPartial("sendercard", senderCardPartial)
	}

	// Correcting the variable name for senderMessagePartial
	senderMessagePartial, err := loadTemplate("Frontend/src/views/partials/sendermessage.hbs")
	if err != nil {
		log.Printf("Warning: Could not load sendermessage partial: %v", err)
	} else {
		raymond.RegisterPartial("sendermessage", senderMessagePartial)
	}

	recipientMessagePartial, err := loadTemplate("Frontend/src/views/partials/recipientmessage.hbs")
	if err != nil {
		log.Printf("Warning: Could not load recipientmessage partial: %v", err)
	} else {
		raymond.RegisterPartial("recipientmessage", recipientMessagePartial)
	}

	filterPartial, err := loadTemplate("Frontend/src/views/partials/filter.hbs")
	if err != nil {
		log.Printf("Warning: Could not load filter partial: %v", err)
	} else {
		raymond.RegisterPartial("filter", filterPartial)
	}

	headerPartial, err := loadTemplate("Frontend/src/views/partials/header.hbs")
	if err != nil {
		log.Printf("Warning: Could not load header partial: %v", err)
	} else {
		raymond.RegisterPartial("header", headerPartial)
	}

	registerPartial, err := loadTemplate("Frontend/src/views/register.hbs")
	if err != nil {
		log.Printf("Warning: Could not load register partial: %v", err)
	} else {
		raymond.RegisterPartial("register", registerPartial)
	}

	errorCardPartial, err := loadTemplate("Frontend/src/views/partials/errorcard.hbs")
	if err != nil {
		log.Printf("Warning: Could not load errorcard partial: %v", err)
	} else {
		raymond.RegisterPartial("errorcard", errorCardPartial)
	}

	log.Println("🚀 Server running on http://0.0.0.0:9081")
	router.Run("0.0.0.0:9081")
}

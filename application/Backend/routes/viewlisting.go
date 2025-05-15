package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"application/Backend/core"
	"application/Backend/database"
	"application/Backend/utils"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterViewListingsRoutes(router *gin.Engine) {
	router.GET("/viewlisting/:id", ProductDetailsMiddleware(), func(c *gin.Context) {
		log.Println("viewlisting: Entering viewlisting route")
		productDetails, exists := c.Get("productDetails")
		if !exists {
			log.Println("viewlisting: Product details not found in context")
			renderErrorPage(c, http.StatusInternalServerError, "Failed to load product details")
			return
		}

		log.Printf("viewlisting: Product details: %+v", productDetails) // Debugging log

		// Load the viewlisting.hbs template
		viewlistingTemplate, err := core.LoadFrontendFile("src/views/viewlisting.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading viewlisting template: %v", err))
			return
		}

		// Render the viewlisting.hbs template with product details
		content, err := raymond.Render(viewlistingTemplate, productDetails)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering viewlisting template: %v", err))
			return
		}

		// Load the layout.hbs template
		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			return
		}

		// Render the layout.hbs template with the viewlisting content
		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "View Listing",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})
}

func RegisterCreateListingRoutes(router *gin.Engine) {
	router.GET("/createlisting", createListingHandler)
	router.POST("/createlisting", submitListingHandler)
}

func createListingHandler(c *gin.Context) {
	createListingTemplate, err := core.LoadFrontendFile("src/views/createlisting.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading create listing template: %v", err))
		return
	}

	content, err := raymond.Render(createListingTemplate, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering create listing template: %v", err))
		return
	}

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Create Listing",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

func submitListingHandler(c *gin.Context) {
	// Handle form submission logic here
	title := c.PostForm("title")
	desc := c.PostForm("description")
	kind := c.PostForm("kind") // "product" or "event"

	// 1) Handle the single uploaded image (use the first if multiple)
	var imageName string
	if fh, err := c.FormFile("images"); err == nil {
		// Ensure the assets dir exists
		os.MkdirAll("assets", os.ModePerm)

		// Save original image
		imageName = filepath.Base(fh.Filename)
		dst := filepath.Join("assets", imageName)
		if err := c.SaveUploadedFile(fh, dst); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Image save error: %v", err))
			return
		}
		// Optionally generate thumbnail now (market.go generates on demand)
		thumbDir := filepath.Join("assets", "thumbnails")
		os.MkdirAll(thumbDir, os.ModePerm)
		_ = utils.GenerateThumbnail(dst, filepath.Join(thumbDir, imageName), 150, 150)
	}

	// 2) Parse price
	priceStr := c.PostForm("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid price")
		return
	}

	sellerID := c.GetInt("user_id")
	category := c.PostForm("category")
	if kind == "product" && category == "" {
		c.String(http.StatusBadRequest, "Category is required for products")
		return
	}

	// 3) Insert into items
	_, err = database.DB.Exec(`
        INSERT INTO items 
          (title, description, price, category, seller_id, image_url, post_date)
        VALUES (?, ?, ?, ?, ?, ?, NOW())
    `, title, desc, price,
		func() string {
			if kind == "product" {
				return category
			}
			return "event"
		}(),
		sellerID,
		imageName,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("DB insert error: %v", err))
		return
	}

	c.String(http.StatusOK, "Listing submitted successfully!")
}

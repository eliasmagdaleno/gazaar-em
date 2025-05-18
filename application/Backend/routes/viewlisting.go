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

		productMap, ok := productDetails.(map[string]interface{})
		if !ok {
			log.Println("viewlisting: productDetails is not a map[string]interface{}")
			renderErrorPage(c, http.StatusInternalServerError, "Internal error")
			return
		}

		if imageURL, ok := productMap["imageURL"].(string); ok {
			productMap["imageURL"] = "../" + imageURL
		}

		log.Printf("viewlisting: Product details: %+v", productMap)

		viewlistingTemplate, err := core.LoadFrontendFile("src/views/viewlisting.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading viewlisting template: %v", err))
			return
		}

		content, err := raymond.Render(viewlistingTemplate, productMap)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering viewlisting template: %v", err))
			return
		}

		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			return
		}

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
	template, err := core.LoadFrontendFile("src/views/createlisting.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading create listing template: %v", err))
		return
	}

	content, err := raymond.Render(template, nil)
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
	title := c.PostForm("title")
	desc := c.PostForm("description")
	kind := c.PostForm("kind")

	var imageName string
	if fh, err := c.FormFile("images"); err == nil {
		os.MkdirAll("assets", os.ModePerm)
		imageName = filepath.Base(fh.Filename)
		dst := filepath.Join("assets", imageName)
		if err := c.SaveUploadedFile(fh, dst); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Image save error: %v", err))
			return
		}
		thumbDir := filepath.Join("assets", "thumbnails")
		os.MkdirAll(thumbDir, os.ModePerm)
		_ = utils.GenerateThumbnail(dst, filepath.Join(thumbDir, imageName), 150, 150)
	}

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
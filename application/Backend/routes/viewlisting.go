package routes

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"application/Backend/core"
	"application/Backend/database"
	"application/Backend/utils"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"

)

func RegisterViewListingsRoutes(router *gin.Engine) {
	router.Use(SignedInMiddleware())
	router.GET("/viewlisting/:id", ProductDetailsMiddleware(), func(c *gin.Context) {
		// log.Println("viewlisting: Entering viewlisting route")

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

	router.POST("/send-message", func(c *gin.Context) {
		productID := c.PostForm("id")
		sellerID := c.PostForm("sellerID")
		message := c.PostForm("message")
		userID := c.GetInt("user_id")
		log.Printf("send-message: userID: %d", userID)


		log.Printf("send-message: productID: %s, sellerID: %s, message: %s", productID, sellerID, message)

		if productID == "" || sellerID == "" || message == "" {
			c.String(http.StatusBadRequest, "Missing required fields")
			return
		}

		roomID, err := strconv.Atoi(productID)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid product ID")
			return
		}
		log.Printf("send-message: roomID: %d", roomID)


		_, err = database.DB.Exec(`
			INSERT INTO Message (sender_id, receiver_id, content, timestamp, room)
			VALUES (?, ?, ?, NOW(), ?)
		`, userID, sellerID, message, roomID)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to save message: %v", err))
			return
		}

		// Redirect back to the product page or show a success message
		c.Redirect(http.StatusSeeOther, "/viewlisting/"+productID)
	})
}

func RegisterCreateListingRoutes(router *gin.Engine) {
	router.GET("/createlisting", createListingHandler)
	router.POST("/createlisting", submitListingHandler)

	router.GET("/selectlocation", selectLocationHandler)
	router.POST("/createlisting/submit", finalizeListingHandler)
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
	priceStr := c.PostForm("price")
	category := c.PostForm("category")
	imageName := ""
	new_base := ""
	if fh, err := c.FormFile("images"); err == nil {
		imageName = filepath.Base(fh.Filename)
		ext := strings.ToLower(filepath.Ext(imageName))
		dst := filepath.Join("frontend/assets/originalImage", imageName) 
		
		if err := c.SaveUploadedFile(fh, dst); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Image save error: %v", err))
			return
		}

		if kind == "product" && category == "" {
			category = "events"
		}


		if kind == "event" {
			category = "events"
		}
		
		thumbDir := filepath.Join("frontend", "assets", "thumbnails")
		thumbPath := filepath.Join(thumbDir, title + ext)

		if err := utils.GenerateThumbnail(dst, thumbPath, 150, 150); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Thumbnail generation error: %v", err))
			return
		}
		new_base = title + ext
	}

	c.SetCookie("listing_title", title, 3600, "/", "", false, true)
	c.SetCookie("listing_description", desc, 3600, "/", "", false, true)
	c.SetCookie("listing_kind", kind, 3600, "/", "", false, true)
	c.SetCookie("listing_image", new_base, 3600, "/", "", false, true)
	c.SetCookie("listing_price", priceStr, 3600, "/", "", false, true)
	c.SetCookie("listing_category", category, 3600, "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/selectlocation")
}

func selectLocationHandler(c *gin.Context) {
	selectLocationTemplate, err := core.LoadFrontendFile("src/views/selectlocation.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading select location template: %v", err)
		return
	}

	content, err := raymond.Render(selectLocationTemplate, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering select location template: %v", err)
		return
	}

	layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading layout template: %v", err)
		return
	}

	output, err := raymond.Render(layoutTemplate, map[string]interface{}{
		"title":   "Select Location",
		"content": raymond.SafeString(content),
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering layout: %v", err)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, output)
}

func finalizeListingHandler(c *gin.Context) {
	sellerID := c.GetInt("user_id")
	location := c.PostForm("location") // From hidden field or button click

	// Read data from cookies
	title, _ := c.Cookie("listing_title")
	desc, _ := c.Cookie("listing_description")
	imageName, _ := c.Cookie("listing_image")
	priceStr, _ := c.Cookie("listing_price")
	category, _ := c.Cookie("listing_category")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		price = 0.0
	}

	_, err = database.DB.Exec(`
        INSERT INTO items 
          (title, description, price, category, seller_id, image_url, post_date, address)
        VALUES (?, ?, ?, ?, ?, ?, NOW(), ?)
    `, title, desc, price, category, sellerID, imageName, location)

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("DB insert error: %v", err))
		return
	}

	// Clear cookies (optional cleanup)
	c.SetCookie("listing_title", "", -1, "/", "", false, true)
	c.SetCookie("listing_description", "", -1, "/", "", false, true)
	c.SetCookie("listing_kind", "", -1, "/", "", false, true)
	c.SetCookie("listing_image", "", -1, "/", "", false, true)
	c.SetCookie("listing_price", "", -1, "/", "", false, true)
	c.SetCookie("listing_category", "", -1, "/", "", false, true)

	c.String(http.StatusOK, "Listing created successfully!")
	c.Redirect(http.StatusSeeOther, "/")
}
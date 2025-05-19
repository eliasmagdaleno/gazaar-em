package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv" // Added import

	"application/Backend/core"
	"application/Backend/database"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
)

func RegisterMessagesRoutes(router *gin.Engine) {
	// GET endpoint to view messages
	router.GET("/messages", func(c *gin.Context) {
		currentRoomIdQuery := c.Query("room")
		userID := c.GetInt("user_id")
		log.Println("messages: User ID from context:", userID)
		if userID == 0 {
			c.String(http.StatusUnauthorized, "User not logged in")
			return
		}
		log.Println("Current user ID:", userID)

		var rows *sql.Rows
		var err error

		if currentRoomIdQuery != "" {
			log.Println("Room ID is not empty, fetching messages for room:", currentRoomIdQuery)
			rows, err = database.DB.Query(`
				SELECT content,
            		DATE_FORMAT(timestamp, '%Y-%m-%d %H:%i') AS ts,
					sender_id,
					receiver_id
          		FROM Message
        		WHERE room = ? 
				AND (sender_id = ? OR receiver_id = ?)
        		ORDER BY timestamp DESC
			`, currentRoomIdQuery, userID, userID)
		} else {
			log.Println("Room ID is empty, fetching all messages for this user")
			rows, err = database.DB.Query(`
				SELECT content,
            		DATE_FORMAT(timestamp, '%Y-%m-%d %H:%i') AS ts,
					sender_id,
					receiver_id
          		FROM Message
        		WHERE sender_id = ? OR receiver_id = ?
         		ORDER BY timestamp DESC
			`, userID, userID)
		}
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("DB error: %v", err))
			return
		}
		defer rows.Close()

		var allMsgs []map[string]interface{}
		for rows.Next() {
			var message, ts string
			var senderID, receiverID int
			if err := rows.Scan(&message, &ts, &senderID, &receiverID); err != nil {
				log.Printf("Error scanning message row: %v", err)
				continue
			}
			allMsgs = append(allMsgs, map[string]interface{}{
				"message":    message,
				"timestamp":  ts,
				"isSender":   senderID == userID,
				"senderID":   senderID,
				"receiverID": receiverID,
			})
		}
		log.Println("All messages for current view:", allMsgs)

		roomRows, err := database.DB.Query(`
			SELECT DISTINCT room FROM Message WHERE sender_id = ? OR receiver_id = ?`, userID, userID)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("DB error fetching rooms: %v", err))
			return
		}
		defer roomRows.Close()

		var rooms []map[string]interface{}
		for roomRows.Next() {
			var rId int // Renamed to avoid conflict with currentRoomIdQuery
			if err := roomRows.Scan(&rId); err != nil {
				log.Printf("Error scanning room ID row: %v", err)
				continue
			}
			var otherUserIDInList int
			err = database.DB.QueryRow(`
				SELECT CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END 
				FROM Message WHERE room = ? AND (sender_id = ? OR receiver_id = ?) 
				LIMIT 1`, userID, rId, userID, userID).Scan(&otherUserIDInList)
			if err != nil {
				log.Printf("Error fetching other user ID for room %d in list: %v", rId, err)
				// Fallback or skip this room from the list if critical info is missing
				otherUserIDInList = 0 // Or some indicator of an issue
			}
			var otherUserNameInList string
			if otherUserIDInList != 0 {
				err = database.DB.QueryRow("SELECT user_name FROM Account WHERE user_id = ?", otherUserIDInList).Scan(&otherUserNameInList)
				if err != nil {
					log.Printf("Error fetching other user name for user ID %d in list: %v", otherUserIDInList, err)
					otherUserNameInList = fmt.Sprintf("User %d", otherUserIDInList)
				}
			} else {
				otherUserNameInList = "Unknown User"
			}
			
			var lastMessage string
			_ = database.DB.QueryRow("SELECT content FROM Message WHERE room = ? ORDER BY timestamp DESC LIMIT 1", rId).Scan(&lastMessage)

			rooms = append(rooms, map[string]interface{}{
				"id":          rId,
				"name":        otherUserNameInList,
				"lastMessage": lastMessage,
			})
		}

		var currentChatOtherUserID int
		var currentChatOtherUserName string

		if currentRoomIdQuery != "" {
			err := database.DB.QueryRow(`
                SELECT 
                    CASE 
                        WHEN m.sender_id = ? THEN m.receiver_id 
                        ELSE m.sender_id 
                    END as other_user_id,
                    a.user_name as other_user_name
                FROM Message m
                JOIN Account a ON a.user_id = (CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END)
                WHERE m.room = ? AND (m.sender_id = ? OR m.receiver_id = ?)
                LIMIT 1`, userID, userID, currentRoomIdQuery, userID, userID).Scan(&currentChatOtherUserID, &currentChatOtherUserName)
			if err != nil {
				if (err == sql.ErrNoRows) {
					log.Printf("No messages in room %s or room does not involve user %d. Cannot determine chat partner.", currentRoomIdQuery, userID)
					currentChatOtherUserName = "Chat" 
					currentChatOtherUserID = 0      
				} else {
					log.Printf("Error fetching other user details for current room %s: %v", currentRoomIdQuery, err)
					currentChatOtherUserName = "Error"
					currentChatOtherUserID = 0
				}
			}
		} else {
			currentChatOtherUserName = "Select a chat"
			currentChatOtherUserID = 0
		}

		messagesTemplate, err := core.LoadFrontendFile("src/views/messages.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading messages template: %v", err))
			return
		}
		layoutTemplate, err := core.LoadFrontendFile("src/views/layouts/layout.hbs")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading layout template: %v", err))
			return
		}
		
		content, err := raymond.Render(messagesTemplate, map[string]interface{}{
			"title":             "View Messages",
			"messages":          allMsgs,
			"rooms":             rooms,
			"roomId":            currentRoomIdQuery,
			"name":              currentChatOtherUserName,    // Name for the header of the active chat
			"receiverIdForForm": currentChatOtherUserID,      // ID for the hidden input in the form
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering messages template: %v", err))
			return
		}

		output, err := raymond.Render(layoutTemplate, map[string]interface{}{
			"title":   "View Messages",
			"content": raymond.SafeString(content),
		})
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering layout: %v", err))
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, output)
	})

	// POST endpoint to send messages
	router.POST("/messages", func(c *gin.Context) {
		message := c.PostForm("message")
		roomIdStr := c.PostForm("room")
		userID := c.GetInt("user_id")
		receiverIdStr := c.PostForm("receiver_id")

		log.Printf("POST /messages: userID: %d, receiverIdStr: %s, message: %s, roomIDStr: %s", userID, receiverIdStr, message, roomIdStr)

		if message == "" || roomIdStr == "" {
			log.Println("Message or RoomID is empty")
			if roomIdStr != "" {
				c.Redirect(http.StatusSeeOther, "/messages?room="+roomIdStr)
			} else {
				c.Redirect(http.StatusSeeOther, "/messages")
			}
			return
		}

		receiverID, err := strconv.Atoi(receiverIdStr)
		if err != nil || receiverID == 0 {
			log.Printf("Invalid receiver_id: '%s'. Error: %v", receiverIdStr, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing receiver ID for the message."})
			// It might be better to redirect with an error message if this is user-facing
			// c.Redirect(http.StatusSeeOther, "/messages?room="+roomIdStr+"&error=invalid_receiver")
			return
		}
		
		// Assuming roomID is also an integer in the database, convert it.
		// If roomID in DB is string, this conversion is not needed for roomID.
		// Based on schema (room INT), it should be an int.
		roomID, err := strconv.Atoi(roomIdStr)
		if err != nil {
			log.Printf("Invalid room_id: '%s'. Error: %v", roomIdStr, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format."})
			return
		}


		_, err = database.DB.Exec(`
			INSERT INTO Message (content, timestamp, sender_id, receiver_id, room)
			VALUES (?, NOW(), ?, ?, ?)
		`, message, userID, receiverID, roomID) // Use parsed receiverID (int) and roomID (int)
		if err != nil {
			log.Printf("Failed to save message to DB: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save message: %v", err)})
			return
		}

		log.Printf("Message sent successfully by user %d to user %d in room %d", userID, receiverID, roomID)
		c.Redirect(http.StatusSeeOther, "/messages?room="+roomIdStr)
	})

	// POST endpoint to delete a chat room
	router.POST("/messages/delete-room", func(c *gin.Context) {
		roomId := c.PostForm("room")
		if roomId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
			return
		}

		// Delete all messages in this room
		_, err := database.DB.Exec("DELETE FROM Message WHERE room = ?", roomId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete messages: %v", err)})
			return
		}

		// Optionally, redirect to /messages (no room selected)
		c.Redirect(http.StatusSeeOther, "/messages")
	})
}

package main

import (
	"application/Backend/database"
	"application/Backend/server"
	"log"
)

func main() {
	// Initialize the database connection
	log.Println("Initializing database...")
	database.InitDB()

	// Start the server
	log.Println("Starting the server...")
	server.StartServer()
}

package main

import (
	"application/Backend/server"
	database "application/Database"
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

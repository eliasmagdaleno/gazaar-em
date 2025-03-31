package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"backend/routes"              // adjust the module path as needed
)

var FrontendDir = os.Getenv("FRONTEND_PATH")
var DB *sql.DB

// LoadFrontendFile loads a file from the front-end directory.
func LoadFrontendFile(filelocation string) (string, error) {
	var sb strings.Builder
	sb.WriteString("../Frontend/")
	sb.WriteString(filelocation)
	abspath, err := filepath.Abs(sb.String())
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(abspath)
	if err != nil {
		return "", err
	}
	// Optionally output file content for debugging.
	// os.Stdout.Write(data)
	return string(data), nil
}

// TestDependency is used to ensure the raymond dependency is kept.
func TestDependency() {
	template := "Hello, {{name}}!"
	result, _ := raymond.Render(template, map[string]string{"name": "Zachary"})
	fmt.Println(result)
}

// initDB initializes the database connection.
func initDB() error {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return fmt.Errorf("DB_DSN environment variable not set")
	}
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return DB.Ping()
}

// StartServer initializes the DB, registers routes, and runs the server.
func StartServer() {
	// Initialize the database.
	if err := initDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	router := gin.Default()
	// Serve static files.
	router.Static("/frontend", "../Frontend")
	router.SetTrustedProxies([]string{"192.168.0.0/24"})

	// Register routes from the routes package.
	routes.RegisterHomeRoutes(router)
	routes.RegisterMemberRoutes(router)
	routes.RegisterVPRoutes(router)
	routes.RegisterSearchRoutes(router)

	log.Println("🚀 Server running on http://0.0.0.0:8081")
	router.Run("0.0.0.0:8081")
}

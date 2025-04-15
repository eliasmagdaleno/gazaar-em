package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dbUser := "remoteuser"
	dbPass := "HorseMomDadHouseThing1!"
	dbHost := "204.236.166.51"
	dbPort := "8081"
	dbName := "SFSU_Marketplace"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MySQL: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Could not ping DB: %v", err)
	}

	log.Println("✅ Successfully connected to MySQL!")
}

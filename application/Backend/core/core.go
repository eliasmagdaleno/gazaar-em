package core

import (
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

//  var DB *sql.DB


// // InitDB sets up a MySQL connection pool using environment variable DB_DSN.
// func InitDB() error {
// 	dsn := os.Getenv("DB_DSN")
// 	if dsn == "" {
// 		return fmt.Errorf("DB_DSN environment variable not set")
// 	}
// 	var err error
// 	DB, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		return err
// 	}
// 	return DB.Ping()
// }


// LoadFrontendFile reads HTML or static file from the frontend directory.
func LoadFrontendFile(filelocation string) (string, error) {
	path := filepath.Join("Frontend", filepath.Clean(filelocation))
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

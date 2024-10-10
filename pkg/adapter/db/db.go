package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	// Get the connection URL from the environment or hardcode for testing
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Error load database url")
	}

	// Open the database connection
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping the database to check the connection
	err = database.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
	return database
}

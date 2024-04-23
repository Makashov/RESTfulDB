package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

// DatabaseConfig holds the configuration for the database connection
type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// ConnectToMySQL creates a connection to the MySQL database and returns the *sql.DB object.
func ConnectToMySQL(cfg DatabaseConfig) (*sql.DB, error) {
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port)

	// Open a new connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open connection to database: %w", err)
	}

	// Verify the connection is successful by pinging the database
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}

func getDatabases(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("could not query databases: %w", err)
	}
	//defer rows.Close()

	var databases []string
	var databaseName string

	for rows.Next() {
		if err := rows.Scan(&databaseName); err != nil {
			log.Fatalf("Error scanning database row: %s", err)
		}
		databases = append(databases, databaseName)
	}

	return databases, nil
}

func main() {
	cfg := DatabaseConfig{
		Username: os.Getenv("MYSQL_ROOT_USER"),
		Password: os.Getenv("MYSQL_ROOT_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
	}

	db, err := ConnectToMySQL(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	databases, err := getDatabases(db)
	if err != nil {
		log.Fatalf("Error querying databases: %s", err)
	}

	fmt.Println("Databases:")
	for _, databaseName := range databases {
		fmt.Println(databaseName)
	}
}

package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/harshgupta9473/webhookDelivery/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	c := config.LoadConfig()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	log.Println("PostgreSQL connected.")
}

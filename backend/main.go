package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/webhookDelivery/db"
	redisHelper "github.com/harshgupta9473/webhookDelivery/redis"
	"github.com/harshgupta9473/webhookDelivery/routes"
	"github.com/harshgupta9473/webhookDelivery/workers/delivery"
	"github.com/harshgupta9473/webhookDelivery/workers/logsCleanup"
	// "github.com/joho/godotenv"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file: ", err)
	// }
	// log.Println("ENV loaded")
	// Initialize database
	db.InitDB()
	redisHelper.InitRedis()
	log.Println("Database initialized")

	// Create tables
	err := db.CreatAllTable()
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
	log.Println("Tables created")

	router := mux.NewRouter()
	routes.RegisterRoutes(router)
	go logsCleanup.RunCleanupTask()
	go delivery.DeliveryWorker()
	headersOk := handlers.AllowedHeaders([]string{"X-event-type", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk)(router)
	s := &http.Server{
		Addr:         ":8080",
		Handler:      corsHandler,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Listening on port :8080")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	sig := <-sigChan
	log.Println("Recieved signal to terminate:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	log.Println("Server exited properly")

}

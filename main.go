package main

import (
	"context"
	"log"
	"time"

	"github.com/deepaksing/FampayAssignment/server"
	"github.com/deepaksing/FampayAssignment/store"
	"github.com/deepaksing/FampayAssignment/store/db/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func fetchAndStoreVideo(ctx context.Context, storeConn *store.Store) {
	ticker := time.NewTicker(40 * time.Second)
	defer ticker.Stop()

	if err := store.FetchAndStore(storeConn); err != nil {
		log.Println("Error fetching and storing videos:", err)
	}

	for {
		select {
		case <-ctx.Done():
			return // Exit goroutine if context is canceled
		case <-ticker.C:
			if err := store.FetchAndStore(storeConn); err != nil {
				log.Println("Error fetching and storing videos:", err)
			}
		}
	}
}

func main() {

	ctx := context.Background()
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create DB connection
	dbConn, err := postgres.NewDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = dbConn.Migrate(ctx)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	storeConn := store.NewStore(dbConn)

	go fetchAndStoreVideo(ctx, storeConn)

	//API routes are registerd here
	server := server.NewServer(storeConn)
	server.StartServer()
}

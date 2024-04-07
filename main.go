package main

import (
	"context"
	"log"

	"github.com/deepaksing/FampayAssignment/server"
	"github.com/deepaksing/FampayAssignment/store"
	"github.com/deepaksing/FampayAssignment/store/db/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
	// Start fetching and storing videos
	store.FetchAndStoreVideos(storeConn)

	server := server.NewServer(storeConn)
	server.StartServer()
}

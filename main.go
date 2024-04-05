package main

import (
	"fmt"
	"log"

	"github.com/deepaksing/FampayAssignment/store/db/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn, err := postgres.NewDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(dbConn)
}

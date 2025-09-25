package main

import (
	"log"

	"e-wallet/config"
	"e-wallet/routes"
)

func main() {
	// Connect to MongoDB
	mongoCfg, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Printf("Successfully connected to database: %s", mongoCfg.DB.Name())

	routes.SetupRoutes(mongoCfg.DB).Run(":8080")
}

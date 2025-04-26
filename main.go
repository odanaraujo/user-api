package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/odanaraujo/user-api/infrastructure/loggers"
	"github.com/odanaraujo/user-api/router"
)

func main() {
	defer loggers.Close()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env found, using environment variables")

	}

	r := router.NewRouter()
	r.Run(":8080")
}

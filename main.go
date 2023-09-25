package main

import (
	"log"
	"os"
	"time"

	"github.com/DamienDuv/rssagg/internal/api"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatalln("DB_URL is not found in the environment")
	}

	server := api.NewServer(dbUrl)

	server.StartScraping(10, time.Minute * 1)

	log.Println("Server starting on port ", portString)
	errSrv := server.StartListening(portString)
	if errSrv != nil {
		log.Fatalln(errSrv)
	}

}

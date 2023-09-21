package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello world")

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatalln("PORT is not found in the environment")
	}

	fmt.Println("PORT:", portString)
}
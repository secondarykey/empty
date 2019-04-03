package main

import (
	"log"
	"os"

	"app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	app.Listen(port)
}

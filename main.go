package main

import (
	"github/kijunpos/internal/app"
	"log"
)

func main() {
	// Initialize the application
	application := app.NewApplication()
	
	// Start the application
	log.Println("Starting KijunPOS application...")
	application.Start()
}

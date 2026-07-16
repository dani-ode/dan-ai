package main

import (
	"log"

	"dan-ai/apps/api/bootstrap"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"portfolio-ai/apps/worker-events/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}

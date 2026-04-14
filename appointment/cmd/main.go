package main

import (
	"log"
	"os"

	"appointment/internal/app"
)

func main() {
	doctorURL := os.Getenv("DOCTOR_URL")
	if doctorURL == "" {
		doctorURL = "http://localhost:8080"
	}

	if err := app.Run(":8081", doctorURL); err != nil {
		log.Fatal(err)
	}
}

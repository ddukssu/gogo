package main

import (
	"log"

	"doc/internal/app"
)

func main() {
	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"diary/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Falha no servidor: %v", err)
	}
}

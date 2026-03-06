package main

import (
	"context"
	"log"

	"github.com/rlapenok/rybakov_test/internal/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

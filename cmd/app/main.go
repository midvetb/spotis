package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/maya-florenko/spotis/internal/app"
	"github.com/maya-florenko/spotis/internal/banner"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	banner.Print()

	if err := app.Init(ctx); err != nil {
		log.Fatal(err)
	}
}

package app

import (
	"context"
	"os"

	"github.com/go-telegram/bot"
)

func Init(ctx context.Context) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(Handler),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_TOKEN"), opts...)
	if err != nil {
		return err
	}

	b.Start(ctx)

	return nil
}

package app

import (
	"context"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Init(ctx context.Context) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, CommandStart),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_TOKEN"), opts...)
	if err != nil {
		return err
	}

	b.Start(ctx)

	return nil
}

func handler(ctx context.Context, b *bot.Bot, u *models.Update) {
	MessageHandler(ctx, b, u)
	InlineHandler(ctx, b, u)
}

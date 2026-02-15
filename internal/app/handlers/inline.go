package handlers

import (
	"context"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/maya-florenko/spotis/internal/deezer"
	"github.com/maya-florenko/spotis/internal/songlink"
	"github.com/maya-florenko/spotis/internal/spotify"
	spot "github.com/zmb3/spotify/v2"
)

func InlineHandler(ctx context.Context, b *bot.Bot, u *models.Update) {
	if u.InlineQuery == nil {
		return
	}

	res, err := spotify.Get(ctx, u.InlineQuery.Query)
	if err != nil {
		return
	}

	id, err := download(ctx, b, u.InlineQuery.Query, res)
	if err != nil {
		return
	}

	results := []models.InlineQueryResult{
		&models.InlineQueryResultCachedAudio{
			ID:          "1",
			AudioFileID: id,
		},
	}

	b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
		InlineQueryID: u.InlineQuery.ID,
		Results:       results,
		CacheTime:     0,
	})
}

func download(ctx context.Context, b *bot.Bot, url string, res *spot.FullTrack) (string, error) {
	u, err := songlink.Link(ctx, url)
	if err != nil {
		return "", err
	}

	file, name, err := deezer.DownloadTrackFromURL(ctx, u)
	if err != nil {
		return "", err
	}

	msg, err := b.SendAudio(ctx, &bot.SendAudioParams{
		ChatID: os.Getenv("TELEGRAM_CHAT_ID"),
		Audio: &models.InputFileUpload{
			Filename: name,
			Data:     file,
		},
		Title: res.Artists[0].Name + " - " + res.Name,
	})
	if err != nil {
		return "", err
	}

	return msg.Audio.FileID, err
}

package main

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/go-telegram/bot"
)

func (app *application) sendServerError(ctx context.Context, b *bot.Bot, chatID int64, err error) {
	sendError(ctx, b, chatID, app.config.locale.Translate("There is a problem in our servers. Please be patient and try later! ⚠️"), err)
}

func sendError(ctx context.Context, b *bot.Bot, chatID int64, msg string, err error) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})

	// Don't log the message if there is no rows in the table
	if errors.Is(err, sql.ErrNoRows) {
		return
	}

	log.Println(err)
}

func contains(slice []int64, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}

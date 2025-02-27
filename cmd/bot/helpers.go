package main

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
)

func sendError(ctx context.Context, b *bot.Bot, chatID int64, msg string, err error) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})

	log.Println(err)
}

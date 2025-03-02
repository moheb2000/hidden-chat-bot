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

func contains(slice []int64, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}

	return false
}

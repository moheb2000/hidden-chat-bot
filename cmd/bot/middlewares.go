package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// addUserInFirstStart checks if a user record is in database or not and if it's not it will add it to database
func (app *application) addUserInFirstStart(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var chatID int64
		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Message.Chat.ID
		} else {
			chatID = update.Message.Chat.ID
		}

		exists, err := app.users.Exists(chatID)
		if err != nil {
			sendError(ctx, b, chatID, "There is a problem in our servers. Please wait a moment and try again...", err)
			return
		}

		// If the user doesn't exist in the database, create a new user record
		if !exists {
			err = app.users.Insert(update.Message.Chat.ID)
			if err != nil {
				sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)
				return
			}
		}

		isOneUser, err := app.users.IsOneUser()
		if err != nil {
			return
		}

		if isOneUser {
			err = app.users.MakeAdmin(chatID)
			if err != nil {
				return
			}
		}

		next(ctx, b, update)
	}
}

// checkIfUserIsBanned is a middleware that runs on every requests to the bot and check if a user is banned or not
func (app *application) checkIfUserIsBanned(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var chatID int64
		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Message.Chat.ID
		} else {
			chatID = update.Message.Chat.ID
		}

		u, err := app.users.GetBychatID(chatID)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Check if a user is ban and not an admin.
		if u.IsBan && !u.IsAdmin {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   app.config.locale.Translate("You are banned by bot's admin and you can't use bot anymore! ⛔❌🔴"),
			})

			return
		}

		next(ctx, b, update)
	}
}

// toggleTypePermission is a middleware that runs before showing premissions page. This allows users to see changes in permissions instantly after clicking on a type inline button
func (app *application) toggleTypePermission(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		pl := strings.Split(update.CallbackQuery.Data, "_")

		// toggle_permission_<per> has two underscores, so the length of callback query data supposed to be 3
		if len(pl) != 3 {
			b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
				ShowAlert:       false,
			})

			return
		}

		per := pl[2]
		app.users.TogglePermission(update.CallbackQuery.Message.Message.Chat.ID, per)

		next(ctx, b, update)
	}
}

func (app *application) settingsChangeLinkAccept(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		// This tells telegram that we are answering the callback query
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		err := app.users.ChangeId(update.CallbackQuery.Message.Message.Chat.ID)
		if err != nil {
			return
		}

		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
		})

		next(ctx, b, update)
	}
}

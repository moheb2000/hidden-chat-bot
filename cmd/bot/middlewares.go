package main

import (
	"context"
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
			app.sendServerError(ctx, b, chatID, err)
			return
		}

		// If the user exists in the database, run the next handler and return
		if exists {
			next(ctx, b, update)
			return
		}

		// All of the code below runs if the user didn't exist in the database for the first time:
		// Insert new user to the database
		err = app.users.Insert(update.Message.Chat.ID)
		if err != nil {
			app.sendServerError(ctx, b, chatID, err)
			return
		}

		isOneUser, err := app.users.IsOneUser()
		if err != nil {
			app.sendServerError(ctx, b, chatID, err)
			return
		}

		if isOneUser {
			err = app.users.MakeAdmin(chatID)
			if err != nil {
				app.sendServerError(ctx, b, chatID, err)
				return
			}
		}

		next(ctx, b, update)

		// Because users can start the bot with deeplink for the first time, they will not have reply buttons till they start the bot again without any deep link, I add this reply button after users first start the bot whether its start with deep link or without
		rkm := models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: app.config.locale.Translate("🔗 Get Hidden Link")},
				},
				{
					{Text: app.config.locale.Translate("⚙️ Settings")},
					{Text: app.config.locale.Translate("ℹ️ About")},
				},
			},
			ResizeKeyboard: true,
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.Message.Chat.ID,
			ReplyMarkup: rkm,
		})
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

		// Becuase checkIfUserIsBanned is the second middleware in chain, if there is an error in getting user by chat_id, it will not be a not found error, Because we check it in addUserInFirstStart middleware to the database.
		u, err := app.users.GetBychatID(chatID)
		if err != nil {
			app.sendServerError(ctx, b, chatID, err)
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
			app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, err)
			return
		}

		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
		})

		next(ctx, b, update)
	}
}

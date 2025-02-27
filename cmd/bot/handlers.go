// Handler functions for bot go in this file
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

// start handles the "/start" command in the bot. If user send this command with an id in the database, user can send hidden message to someone with that id in the database and if this command doesn't have an id, welcome screen will be shown to the user
func (app *application) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	// TODO: check if this code must go to a middleware to run on all requests or not
	// Check if the user exists in the database records or not and if it doesn't create a new user record
	exists, err := app.users.Exists(update.Message.Chat.ID)
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)
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

	// ml is a variable to handle deep links sended to the bot for sending messages
	ml := strings.Split(update.Message.Text, " ")
	// Check if there is an id in the command or not
	if len(ml) == 2 {
		// run sendState to change the state of the user
		app.sendState(ctx, b, update.Message, ml[1])

		// This is the end of code runs when there is a deeplink with and id for sending messages
		return
	}

	// Modify this section to show the keyboard when starting the app no matter what command used
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello, Welcome to this app!",
		ReplyMarkup: models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "Create hidden link"},
				},
				{
					{Text: "Settings"},
					{Text: "About"},
				},
			},
			ResizeKeyboard: true,
		},
	})
}

// getHiddenLink sends the anonymous link when the user clicks on this button
func (app *application) getHiddenLink(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Get the user by chat id in the telegram
	u, err := app.users.GetBychatID(update.Message.Chat.ID)
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("t.me/hidden_chat_moheb2000_bot?start=%s", u.ID.String()),
	})
}

// about runs when the user clicks on the about reply command and will show the about text to the user
func about(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "This is the about page",
	})
}

// send handler runs on every message users send to the bot. It will send the message to another user or show an error base on the sending state of the user in the database
func (app *application) send(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Get user from database
	u, err := app.users.GetBychatID(update.Message.Chat.ID)
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)
		return
	}

	// Check if the user is in sending mode or not
	if !u.IsSending {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "You are not in sending mode...",
		})

		return
	}

	// Get the recipient user from the database
	ru, err := app.users.Get(u.RecipientID)
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)

		return
	}

	// Send the message with a inline button to the recipient user. If user clicks on reply button, a callback query with reply_<recipient_id> data will be send.

	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Reply", CallbackData: "reply_" + u.ID.String()},
			},
		},
	}

	// Check if the message type is text
	if update.Message.Text != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      ru.ChatID,
			Text:        update.Message.Text,
			ReplyMarkup: ibm,
		})
		// Check if the message type is sticker
	} else if update.Message.Sticker != nil {
		b.SendSticker(ctx, &bot.SendStickerParams{
			ChatID:      ru.ChatID,
			Sticker:     &models.InputFileString{Data: update.Message.Sticker.FileID},
			ReplyMarkup: ibm,
		})
		// Check if the message type is a GIF
	} else if update.Message.Animation != nil {
		b.SendAnimation(ctx, &bot.SendAnimationParams{
			ChatID:      ru.ChatID,
			Animation:   &models.InputFileString{Data: update.Message.Animation.FileID},
			Caption:     update.Message.Caption,
			ReplyMarkup: ibm,
		})
		// Check if the message type is photo
	} else if update.Message.Photo != nil {
		b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:      ru.ChatID,
			Photo:       &models.InputFileString{Data: update.Message.Photo[len(update.Message.Photo)-1].FileID},
			Caption:     update.Message.Caption,
			ReplyMarkup: ibm,
		})
		// Check if the message type is video
	} else if update.Message.Video != nil {
		b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:      ru.ChatID,
			Video:       &models.InputFileString{Data: update.Message.Video.FileID},
			Caption:     update.Message.Caption,
			ReplyMarkup: ibm,
		})
		// Check if the message type is voice
	} else if update.Message.Voice != nil {
		b.SendVoice(ctx, &bot.SendVoiceParams{
			ChatID:      ru.ChatID,
			Voice:       &models.InputFileString{Data: update.Message.Voice.FileID},
			ReplyMarkup: ibm,
		})
		// Check if the message type is audio
	} else if update.Message.Audio != nil {
		b.SendAudio(ctx, &bot.SendAudioParams{
			ChatID:      ru.ChatID,
			Audio:       &models.InputFileString{Data: update.Message.Audio.FileID},
			Caption:     update.Message.Caption,
			ReplyMarkup: ibm,
		})
		// Check if the message type is document
	} else if update.Message.Document != nil {
		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:      ru.ChatID,
			Document:    &models.InputFileString{Data: update.Message.Document.FileID},
			ReplyMarkup: ibm,
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Your message type doesn't supported with this bot.",
		})

		// We return here because we don't want to leave the sending state just becuase the message type is not supported by the bot
		return
	}

	// After sending the message, user will leave the sending state
	app.users.LeaveSendingState(u.ChatID)
}

func (app *application) reply(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	ql := strings.Split(update.CallbackQuery.Data, "_")
	if len(ql) != 2 {
		sendError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, "A problem happened with data sended to server. Please try again later", nil)
		return
	}

	app.sendState(ctx, b, update.CallbackQuery.Message.Message, ql[1])
}

func (app *application) sendState(ctx context.Context, b *bot.Bot, message *models.Message, rIDString string) {
	// recipientID is the second parameter of deep link
	recipientID, err := uuid.Parse(rIDString)
	// Check if the link is valid
	if err != nil {
		sendError(ctx, b, message.Chat.ID, "This link is not valid. Maybe you need to contact somehow to the link's owner and tell this problem", err)
		return
	}

	// Check if the link is valid
	_, err = app.users.Get(recipientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(ctx, b, message.Chat.ID, "This link is not valid. Maybe you need to contact somehow to the link's owner and tell this problem", err)
		} else {
			sendError(ctx, b, message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)
		}

		return
	}

	// If all of checks pass, enter the sending state and send a message to the user
	app.users.EnterSendingState(message.Chat.ID, recipientID)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: message.Chat.ID,
		Text:   "You can now send an anonymous message:",
	})
}

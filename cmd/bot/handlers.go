// Handler functions for bot go in this file
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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
		ChatID:    update.Message.Chat.ID,
		Text:      app.config.locale.Translate("start_message"),
		ParseMode: models.ParseModeMarkdownV1,
		ReplyMarkup: models.ReplyKeyboardMarkup{
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

// settings runs when the user clicks on the settings reply command and will show a message with inline buttons for managing user preferences
func (app *application) settings(ctx context.Context, b *bot.Bot, update *models.Update) {
	app.showSettings(ctx, b, update, false)
}

func (app *application) backToSettings(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This tells telegram that we are answering the callback query
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	app.showSettings(ctx, b, update, true)
}

func (app *application) showSettings(ctx context.Context, b *bot.Bot, update *models.Update, edit bool) {
	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: app.config.locale.Translate("🚫 Message Restrictions"), CallbackData: "settings_allowed_types"},
			},
		},
	}

	tm := app.config.locale.Translate("settings_message")

	if !edit {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        tm,
			ReplyMarkup: ibm,
		})
	} else {
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        tm,
			ReplyMarkup: ibm,
		})
	}
}

// settingsAllowedTypes runs when user clicks on set allowed message types inline button in settings page and change the settings message to the corresponding text
func (app *application) settingsAllowedTypes(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This tells telegram that we are answering the callback query
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	// Get allowed types from database
	at, err := app.users.GetAllowedTypes(update.CallbackQuery.Message.Message.Chat.ID)
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)

		return
	}

	// ibm is an empty inline button message that we append buttons later
	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}

	// Loop over types of messages and set them allowed or disallowed based on database response
	// TODO: I will change some logic here, so every two inline button be in the same row for more compact response.
	typesBool := []bool{at.Text, at.Sticker, at.Gif, at.Photo, at.Video, at.Voice, at.Audio, at.Document}
	typesString := []string{
		app.config.locale.Translate("Text"),
		app.config.locale.Translate("Sticker"),
		app.config.locale.Translate("Gif"),
		app.config.locale.Translate("Photo"),
		app.config.locale.Translate("Video"),
		app.config.locale.Translate("Voice"),
		app.config.locale.Translate("Audio"),
		app.config.locale.Translate("Document"),
	}
	for i := range typesBool {
		r := ""
		if typesBool[i] {
			r += "✅ "
		} else {
			r += "❌ "
		}
		r += typesString[i]
		ibm.InlineKeyboard = append(ibm.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         r,
				CallbackData: "toggle_permission_" + strings.ToLower(typesString[i]),
			},
		})
	}

	// Add a button to back to the main settings message
	ibm.InlineKeyboard = append(ibm.InlineKeyboard, []models.InlineKeyboardButton{
		{
			Text:         app.config.locale.Translate("⬅️ Back"),
			CallbackData: "back_to_settings",
		},
	})

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        app.config.locale.Translate("message_restrictions_settings_message"),
		ReplyMarkup: ibm,
	})
}

// about runs when the user clicks on the about reply command and will show the about text to the user
func (app *application) about(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      app.config.locale.Translate("about_message"),
		ParseMode: models.ParseModeMarkdownV1,
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
			Text:   app.config.locale.Translate("You are not in sending mode... ⛔"),
		})

		return
	}

	// Get the recipient user from the database
	ru, err := app.users.Get(u.RecipientID)
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)

		return
	}

	// Check if recipient user blocks sender or not
	if contains(ru.Blocks, update.Message.Chat.ID) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   app.config.locale.Translate("You blocked by the user you're trying to send the message! 🔒😿"),
		})

		return
	}

	at, err := app.users.GetAllowedTypes(ru.ChatID)
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)

		return
	}

	// Send the message with a inline button to the recipient user. If user clicks on reply button, a callback query with reply_<recipient_id> data will be send.
	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: app.config.locale.Translate("💬 Reply"), CallbackData: "reply_" + u.ID.String()},
			},
			{
				{Text: app.config.locale.Translate("🔒 Block"), CallbackData: "block_" + strconv.FormatInt(u.ChatID, 10)},
				{Text: app.config.locale.Translate("🚨 Report"), CallbackData: "report"},
			},
			{
				// TODO: Add changing ad link in admin panel after implementing admin user
				{Text: "Ad", URL: "https://t.me/Otazy"},
			},
		},
	}

	// Check if the message type is text and it's an allowed type
	if update.Message.Text != "" && at.Text {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      ru.ChatID,
			Text:        update.Message.Text,
			ReplyMarkup: ibm,
		})
		// Check if the message type is sticker and it's an allowed type
	} else if update.Message.Sticker != nil && at.Sticker {
		b.SendSticker(ctx, &bot.SendStickerParams{
			ChatID:      ru.ChatID,
			Sticker:     &models.InputFileString{Data: update.Message.Sticker.FileID},
			ReplyMarkup: ibm,
		})
		// Check if the message type is a GIF and it's an allowed type
	} else if update.Message.Animation != nil && at.Gif {
		b.SendAnimation(ctx, &bot.SendAnimationParams{
			ChatID:      ru.ChatID,
			Animation:   &models.InputFileString{Data: update.Message.Animation.FileID},
			Caption:     update.Message.Caption,
			ReplyMarkup: ibm,
		})
		// Check if the message type is photo and it's an allowed type
	} else if update.Message.Photo != nil && at.Photo {
		b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:      ru.ChatID,
			Photo:       &models.InputFileString{Data: update.Message.Photo[len(update.Message.Photo)-1].FileID},
			Caption:     update.Message.Caption,
			ReplyMarkup: ibm,
		})
		// Check if the message type is video and it's an allowed type
	} else if update.Message.Video != nil && at.Video {
		b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:      ru.ChatID,
			Video:       &models.InputFileString{Data: update.Message.Video.FileID},
			Caption:     update.Message.Caption,
			ReplyMarkup: ibm,
		})
		// Check if the message type is voice and it's an allowed type
	} else if update.Message.Voice != nil && at.Voice {
		b.SendVoice(ctx, &bot.SendVoiceParams{
			ChatID:      ru.ChatID,
			Voice:       &models.InputFileString{Data: update.Message.Voice.FileID},
			ReplyMarkup: ibm,
		})
		// Check if the message type is audio and it's an allowed type
	} else if update.Message.Audio != nil && at.Audio {
		b.SendAudio(ctx, &bot.SendAudioParams{
			ChatID:      ru.ChatID,
			Audio:       &models.InputFileString{Data: update.Message.Audio.FileID},
			Caption:     update.Message.Caption,
			ReplyMarkup: ibm,
		})
		// Check if the message type is document and it's an allowed type
	} else if update.Message.Document != nil && at.Document {
		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:      ru.ChatID,
			Document:    &models.InputFileString{Data: update.Message.Document.FileID},
			ReplyMarkup: ibm,
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   app.config.locale.Translate("Your message type is limited by reciever or isn't supported by this bot. 🔒🥹"),
		})

		// We return here because we don't want to leave the sending state just becuase the message type is not supported by the bot
		return
	}

	// After sending the message, user will leave the sending state
	app.users.LeaveSendingState(u.ChatID)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   app.config.locale.Translate("Your message sended successfully! 📨😍"),
	})
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
		Text:   app.config.locale.Translate("Send your message ✏️:"),
	})
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

// report runs when a user wants to report a message to admin
// TODO: Add report functionallity after adding admin feature
func report(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}

// block runs when a user clicks on block inline button below a message to blocks the sender message
func (app *application) block(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	bl := strings.Split(update.CallbackQuery.Data, "_")

	if len(bl) != 2 {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", nil)
		return
	}

	// block is a handler function to answers to callback query data which is started by block_ or unblock_. This line of code checks if we are trying to block someone or unblock
	isBlocking := false
	if bl[0] == "block" {
		isBlocking = true
	}

	// Convert chat id from string to int64
	blockChatID, err := strconv.ParseInt(bl[1], 10, 64)
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)
		return
	}

	// Based on isBlocking boolean value, we decide to add or remove chat id of sender to blocks array in database or not
	if isBlocking {
		err = app.users.AddBlockArray(update.CallbackQuery.Message.Message.Chat.ID, blockChatID)
	} else {
		err = app.users.RemoveBlockArray(update.CallbackQuery.Message.Message.Chat.ID, blockChatID)
	}
	if err != nil {
		sendError(ctx, b, update.Message.Chat.ID, "There is a problem in our servers. Please wait a moment and try again...", err)
		return
	}

	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}

	// Retrieve the inline buttons from the callback query
	ibm.InlineKeyboard = update.CallbackQuery.Message.Message.ReplyMarkup.InlineKeyboard
	// Change the block button based on block/unblock users
	// TODO: This is hardcoded logic to change the inline button for blocking feature, but need to be replaced with a more generic approach
	if isBlocking {
		ibm.InlineKeyboard[1][0].Text = app.config.locale.Translate("🔓 Unblock")
		ibm.InlineKeyboard[1][0].CallbackData = "unblock_" + strings.Split(ibm.InlineKeyboard[1][0].CallbackData, "_")[1]
	} else {
		ibm.InlineKeyboard[1][0].Text = app.config.locale.Translate("🔒 Block")
		ibm.InlineKeyboard[1][0].CallbackData = "block_" + strings.Split(ibm.InlineKeyboard[1][0].CallbackData, "_")[1]
	}

	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: ibm,
	})

	t := ""
	if isBlocking {
		t = app.config.locale.Translate("Blocked succesfully! 🔒✅")
	} else {
		t = app.config.locale.Translate("Unblocked successfully! 🔓✅")
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   t,
	})
}

// Handler functions for bot go in this file
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

// start handles the "/start" command in the bot. If user send this command with an id in the database, user can send hidden message to someone with that id in the database and if this command doesn't have an id, welcome screen will be shown to the user
func (app *application) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	// ml is a variable to handle deep links sended to the bot for sending messages
	ml := strings.Split(update.Message.Text, " ")
	// Check if there is an id in the command or not
	if len(ml) == 2 {
		// run sendState to change the state of the user
		app.sendState(ctx, b, update.Message, ml[1])

		// This is the end of code runs when there is a deeplink with and id for sending messages
		return
	}

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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        app.config.locale.Translate("start_message", app.config.name),
		ParseMode:   models.ParseModeMarkdownV1,
		ReplyMarkup: rkm,
	})
}

// getHiddenLink sends the anonymous link when the user clicks on this button or when change hidden link callback query is sended
func (app *application) getHiddenLink(ctx context.Context, b *bot.Bot, update *models.Update) {
	var chatID int64
	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Message.Chat.ID
	} else {
		chatID = update.Message.Chat.ID
	}
	// Get the user by chat id in the telegram
	u, err := app.users.GetBychatID(chatID)
	if err != nil {
		app.sendServerError(ctx, b, chatID, err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("t.me/hidden_chat_moheb2000_bot?start=%s", u.ID.String()),
	})
}

// settings runs when the user clicks on the settings reply command and will show a message with inline buttons for managing user preferences
func (app *application) settings(ctx context.Context, b *bot.Bot, update *models.Update) {
	app.showSettings(ctx, b, update, false)
}

// backToSettings runs when the user clicks on back button in settings menu
func (app *application) backToSettings(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This tells telegram that we are answering the callback query
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	app.showSettings(ctx, b, update, true)
}

// showSettings is a helper function that shows the settings message
func (app *application) showSettings(ctx context.Context, b *bot.Bot, update *models.Update, edit bool) {
	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: app.config.locale.Translate("🔗 Change Hidden Link"), CallbackData: "settings_change_hidden_link"},
			},
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

// settingsChangeLink shows a confirmation message to see if the user really wants to change the hidden chat link or not
func (app *application) settingsChangeLink(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This tells telegram that we are answering the callback query
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: app.config.locale.Translate("Yes, I'm sure ✅"), CallbackData: "settings_change_hidden_link_accept"},
				{Text: app.config.locale.Translate("Cancel ❌"), CallbackData: "settings_change_hidden_link_cancel"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        app.config.locale.Translate("Are you sure you want change the hidden link? Your previous link will be invalid! 🚨🚫"),
		ReplyMarkup: ibm,
	})
}

// settingsChangeLinkCancel runs when user clicks on cancel changing hidden chat link inline button and deleted the confirmation message
func (app *application) settingsChangeLinkCancel(ctx context.Context, b *bot.Bot, update *models.Update) {
	// This tells telegram that we are answering the callback query
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
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
		app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, err)
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
		app.sendServerError(ctx, b, update.Message.Chat.ID, err)
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
		if errors.Is(err, sql.ErrNoRows) {
			sendError(ctx, b, update.Message.Chat.ID, app.config.locale.Translate("The user you trying to send this message to, changes the hidden link. Maybe you could ask the user to get the new one! ⚠️"), err)
			app.users.LeaveSendingState(u.ChatID)
			return
		}

		app.sendServerError(ctx, b, update.Message.Chat.ID, err)
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
		app.sendServerError(ctx, b, update.Message.Chat.ID, err)

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
				{Text: app.config.locale.Translate("🚨 Report"), CallbackData: "report_" + strconv.FormatInt(u.ChatID, 10)},
			},
		},
	}

	if os.Getenv("AD_TEXT") != "" && os.Getenv("AD_URL") != "" {
		ibm.InlineKeyboard = append(ibm.InlineKeyboard, []models.InlineKeyboardButton{
			{Text: os.Getenv("AD_TEXT"), URL: os.Getenv("AD_URL")},
		})
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
		sendError(ctx, b, message.Chat.ID, app.config.locale.Translate("This link is not valid. Maybe you need to contact somehow to the link's owner and tell this problem! ⚠️"), err)
		return
	}

	// Check if the link is valid
	_, err = app.users.Get(recipientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(ctx, b, message.Chat.ID, app.config.locale.Translate("This link is not valid. Maybe you need to contact somehow to the link's owner and tell this problem! ⚠️"), err)
		} else {
			app.sendServerError(ctx, b, message.Chat.ID, err)
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

// report runs when a user wants to report a message to admin
func (app *application) report(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	rl := strings.Split(update.CallbackQuery.Data, "_")
	if len(rl) != 2 {
		return
	}

	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: app.config.locale.Translate("🔴 Ban"), CallbackData: "ban_" + rl[1]},
			},
		},
	}

	ad, err := app.users.GetAdmin()
	if err != nil {
		app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, err)
		return
	}

	b.CopyMessage(ctx, &bot.CopyMessageParams{
		ChatID:      ad.ChatID,
		FromChatID:  update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: ibm,
	})
}

// ban runs when admin clicks on ban inline button below a message to ban the user from using bot
// TODO: User default telegram api for banning instead of a custom logic
func (app *application) ban(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	bl := strings.Split(update.CallbackQuery.Data, "_")

	if len(bl) != 2 {
		app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, nil)
		return
	}

	// ban is a handler function to answers to callback query data which is started by ban_ or unban_. This line of code checks if we are trying to ban someone or unban
	isBanning := false
	if bl[0] == "ban" {
		isBanning = true
	}

	// Convert chat id from string to int64
	banChatID, err := strconv.ParseInt(bl[1], 10, 64)
	if err != nil {
		app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, err)
		return
	}

	// Based on isBanning boolean value, we decide to change the value of is_ban field for the corresponding user to true or false
	if isBanning {
		err = app.users.Ban(banChatID)
	} else {
		err = app.users.Unban(banChatID)
	}
	if err != nil {
		app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, err)
		return
	}

	ibm := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}

	// Change the ban button based on ban/unban users
	if isBanning {
		ibm.InlineKeyboard = append(ibm.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         app.config.locale.Translate("🟢 Unban"),
				CallbackData: "unban_" + bl[1],
			},
		})
	} else {
		ibm.InlineKeyboard = append(ibm.InlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         app.config.locale.Translate("🔴 Ban"),
				CallbackData: "ban_" + bl[1],
			},
		})
	}

	b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: ibm,
	})

	t := ""
	if isBanning {
		t = app.config.locale.Translate("Banned succesfully! 🔴✅")
	} else {
		t = app.config.locale.Translate("Unbanned successfully! 🟢✅")
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   t,
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
		app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, nil)
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
		app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, err)
		return
	}

	// Based on isBlocking boolean value, we decide to add or remove chat id of sender to blocks array in database or not
	if isBlocking {
		err = app.users.AddBlockArray(update.CallbackQuery.Message.Message.Chat.ID, blockChatID)
	} else {
		err = app.users.RemoveBlockArray(update.CallbackQuery.Message.Message.Chat.ID, blockChatID)
	}
	if err != nil {
		app.sendServerError(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, err)
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

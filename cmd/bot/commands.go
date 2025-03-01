package main

import (
	"github.com/go-telegram/bot"
)

// newBot creates a new bot with some options, add handlers to it and return the newly created bot and an error if available
func (app *application) newBot(token string) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(app.send),
		bot.WithCallbackQueryDataHandler("reply_", bot.MatchTypePrefix, app.reply),
		// bot.WithCallbackQueryDataHandler("toggle_permission_", bot.MatchTypePrefix, app.toggleTypePermission),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	// Register handlers
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, app.start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "About", bot.MatchTypeExact, about)
	b.RegisterHandler(bot.HandlerTypeMessageText, "Settings", bot.MatchTypeExact, settings)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "settings_allowed_types", bot.MatchTypeExact, app.settingsAllowedTypes)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "toggle_permission_", bot.MatchTypePrefix, app.settingsAllowedTypes, app.toggleTypePermission)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "back_to_settings", bot.MatchTypeExact, backToSettings)
	b.RegisterHandler(bot.HandlerTypeMessageText, "Create hidden link", bot.MatchTypeExact, app.getHiddenLink)

	return b, nil
}

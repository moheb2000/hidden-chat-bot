package main

import (
	"github.com/go-telegram/bot"
)

// newBot creates a new bot with some options, add handlers to it and return the newly created bot and an error if available
func (app *application) newBot(token string) (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithMiddlewares(app.addUserInFirstStart, app.checkIfUserIsBanned),
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
	b.RegisterHandler(bot.HandlerTypeMessageText, app.config.locale.Translate("ℹ️ About"), bot.MatchTypeExact, app.about)
	b.RegisterHandler(bot.HandlerTypeMessageText, app.config.locale.Translate("⚙️ Settings"), bot.MatchTypeExact, app.settings)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "settings_change_hidden_link", bot.MatchTypeExact, app.settingsChangeLink)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "settings_change_hidden_link_accept", bot.MatchTypeExact, app.getHiddenLink, app.settingsChangeLinkAccept)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "settings_change_hidden_link_cancel", bot.MatchTypeExact, app.settingsChangeLinkCancel)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "settings_allowed_types", bot.MatchTypeExact, app.settingsAllowedTypes)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "toggle_permission_", bot.MatchTypePrefix, app.settingsAllowedTypes, app.toggleTypePermission)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "back_to_settings", bot.MatchTypeExact, app.backToSettings)
	b.RegisterHandler(bot.HandlerTypeMessageText, app.config.locale.Translate("🔗 Get Hidden Link"), bot.MatchTypeExact, app.getHiddenLink)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "block_", bot.MatchTypePrefix, app.block)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "unblock_", bot.MatchTypePrefix, app.block)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "report_", bot.MatchTypePrefix, app.report)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "ban_", bot.MatchTypePrefix, app.ban)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "unban_", bot.MatchTypePrefix, app.ban)
	return b, nil
}

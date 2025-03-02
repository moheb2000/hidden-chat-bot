package localizer

import (
	_ "github.com/moheb2000/hidden-chat-bot/internal/translations"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Localizer struct {
	ID      string
	printer *message.Printer
}

var locales = []Localizer{
	{
		// English
		ID:      "en-us",
		printer: message.NewPrinter(language.MustParse("en-us")),
	},
	{
		// Persian
		ID:      "fa-ir",
		printer: message.NewPrinter(language.MustParse("fa-ir")),
	},
}

func Get(id string) (Localizer, bool) {
	for _, locale := range locales {
		if id == locale.ID {
			return locale, true
		}
	}

	return Localizer{}, false
}

func (l Localizer) Translate(key message.Reference, args ...interface{}) string {
	return l.printer.Sprintf(key, args...)
}

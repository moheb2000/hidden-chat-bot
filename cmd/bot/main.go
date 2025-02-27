// Hidden Chat Bot is a telegram bot that lets users recieve anonymous messages without showing their ids to another one.
package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/moheb2000/hidden-chat-bot/internal/models"
)

type application struct {
	users *models.UserModel
}

// main creates a new bot and starts it
func main() {
	// name := flag.String("name", "Hidden Chat Bot", "The name of the bot")
	// flag.Parse()

	// creates a new pool connection with postgres. If an error occures it will end the bot with log.Fatal()
	db, err := openDB(os.Getenv("HIDDEN_CHAT_DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database connection pool established")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Creating a new application instance for dependency injection to other functions
	app := &application{
		users: &models.UserModel{DB: db},
	}

	// newBot creates a new bot, add handlers to it in commands.go file and then returns the newly created bot
	b, err := app.newBot(os.Getenv("HIDDEN_CHAT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Hidden Chat Bot...")
	b.Start(ctx)
}

// openDB creates a new pool connection with postgres database
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

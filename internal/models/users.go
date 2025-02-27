package models

import (
	"database/sql"

	"github.com/google/uuid"
)

// The User type is a struct for storing each user data
type User struct {
	ID          uuid.UUID
	ChatID      int64
	IsSending   bool
	RecipientID uuid.UUID
}

// The UserModel type has access to an instance of application's database
type UserModel struct {
	DB *sql.DB
}

// Insert adds a user record for each chatID. Each user will have a id with uuid type that will be used for anonymous messaging between users instead of usernames
func (m *UserModel) Insert(chatID int64) error {
	stmt := "INSERT INTO users (chat_id) VALUES ($1)"

	_, err := m.DB.Exec(stmt, chatID)
	if err != nil {
		return err
	}

	return nil
}

// Exists checks if a user exists or not in the database records and returns a boolean value
func (m *UserModel) Exists(chatID int64) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT 1 FROM users WHERE chat_id=$1)"

	err := m.DB.QueryRow(stmt, chatID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Get returns a user with the uuid token specified
func (m *UserModel) Get(id uuid.UUID) (*User, error) {
	stmt := `SELECT id, chat_id, is_sending, recipient_id FROM users
	WHERE id = $1`

	u := &User{}

	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.ChatID, &u.IsSending, &u.RecipientID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Get returns a user with the chat id in the telegram specified
func (m *UserModel) GetBychatID(chatID int64) (*User, error) {
	stmt := `SELECT id, chat_id, is_sending, recipient_id FROM users
	WHERE chat_id = $1`

	u := &User{}

	err := m.DB.QueryRow(stmt, chatID).Scan(&u.ID, &u.ChatID, &u.IsSending, &u.RecipientID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// EnterSendingState changes the is_sending and id_recipient in the database, so when user sends a message it will interpret as a message to be send to someone else
func (m *UserModel) EnterSendingState(chatID int64, recipientID uuid.UUID) error {
	return m.updateSendingState(chatID, true, recipientID)
}

// LeaveSendingState changes the is_sending and id_recipient in the database to default values, so messages sended by users doesn't interpret as sending messages
func (m *UserModel) LeaveSendingState(chatID int64) error {
	return m.updateSendingState(chatID, false, uuid.Nil)
}

// updateSendingState is a helper function and update is_sending and recipient_id to the arguments provided by other functions
func (m *UserModel) updateSendingState(chatID int64, isSending bool, recipientID uuid.UUID) error {
	stmt := `UPDATE users SET is_sending = $1, recipient_id = $2
	WHERE chat_id = $3`

	_, err := m.DB.Exec(stmt, isSending, recipientID, chatID)
	if err != nil {
		return err
	}

	return nil
}

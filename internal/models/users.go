package models

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AllowedTypes struct {
	Text     bool
	Sticker  bool
	Gif      bool
	Photo    bool
	Video    bool
	Voice    bool
	Audio    bool
	Document bool
}

// The User type is a struct for storing each user data
type User struct {
	ID          uuid.UUID
	ChatID      int64
	IsSending   bool
	RecipientID uuid.UUID
	Blocks      []int64
	IsBan       bool
	IsAdmin     bool
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
	stmt := `SELECT id, chat_id, is_sending, recipient_id, blocks, is_ban, is_admin FROM users
	WHERE id = $1`

	u := &User{}

	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.ChatID, &u.IsSending, &u.RecipientID, pq.Array(&u.Blocks), &u.IsBan, &u.IsAdmin)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Get returns a user with the chat id in the telegram specified
func (m *UserModel) GetBychatID(chatID int64) (*User, error) {
	stmt := `SELECT id, chat_id, is_sending, recipient_id, blocks, is_ban, is_admin FROM users
	WHERE chat_id = $1`

	u := &User{}

	err := m.DB.QueryRow(stmt, chatID).Scan(&u.ID, &u.ChatID, &u.IsSending, &u.RecipientID, pq.Array(&u.Blocks), &u.IsBan, &u.IsAdmin)
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

// GetAllowedTypes returns a struct that have all permission types for a single user.
func (m *UserModel) GetAllowedTypes(chatID int64) (*AllowedTypes, error) {
	stmt := `SELECT allow_text, allow_sticker, allow_gif, allow_photo, allow_video, allow_voice, allow_audio, allow_document FROM users
	WHERE chat_id = $1`

	at := &AllowedTypes{}

	err := m.DB.QueryRow(stmt, chatID).Scan(&at.Text, &at.Sticker, &at.Gif, &at.Photo, &at.Video, &at.Voice, &at.Audio, &at.Document)
	if err != nil {
		return nil, err
	}

	return at, nil
}

func (m *UserModel) TogglePermission(chatID int64, per string) error {
	at, err := m.GetAllowedTypes(chatID)
	if err != nil {
		return err
	}

	stmt := "UPDATE users SET allow_" + per + " = $1 WHERE chat_id = $2"
	up := false
	switch per {
	case "text":
		up = !at.Text
	case "sticker":
		up = !at.Sticker
	case "gif":
		up = !at.Gif
	case "photo":
		up = !at.Photo
	case "video":
		up = !at.Video
	case "voice":
		up = !at.Voice
	case "audio":
		up = !at.Audio
	case "document":
		up = !at.Document
	}

	_, err = m.DB.Exec(stmt, up, chatID)
	if err != nil {
		return err
	}

	return nil
}

// AddBlockArray adds a chat id to the blocks array in database
func (m *UserModel) AddBlockArray(chatID int64, blockChatID int64) error {
	stmt := `UPDATE users SET blocks = array_append(blocks, $1)
	WHERE chat_id = $2`

	_, err := m.DB.Exec(stmt, blockChatID, chatID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveBlockArray removes a chat id to the blocks array in database
func (m *UserModel) RemoveBlockArray(chatID int64, blockChatID int64) error {
	stmt := `UPDATE users SET blocks = array_remove(blocks, $1)
	WHERE chat_id = $2`

	_, err := m.DB.Exec(stmt, blockChatID, chatID)
	if err != nil {
		return err
	}

	return nil
}

// MakeAdmin updates a user record and change is_admin to true
func (m *UserModel) MakeAdmin(chatID int64) error {
	stmt := `UPDATE users SET is_admin = TRUE
	WHERE chat_id = $1`

	_, err := m.DB.Exec(stmt, chatID)
	if err != nil {
		return err
	}

	return nil
}

// IsOneUser checks if there is only one record in the users table or not. This function is used in confunction with MakeAdmin to change is_admin field for the first user of the bot to true
func (m *UserModel) IsOneUser() (bool, error) {
	var isOneUser bool
	stmt := "SELECT (COUNT(*) = 1) FROM users"
	err := m.DB.QueryRow(stmt).Scan(&isOneUser)
	if err != nil {
		return false, nil
	}

	return isOneUser, nil
}

// GetAdmin returns the first user in database that has is_admin=true
func (m *UserModel) GetAdmin() (*User, error) {
	stmt := `SELECT id, chat_id, is_sending, recipient_id, blocks, is_ban, is_admin FROM users
	WHERE is_admin = TRUE ORDER BY serial LIMIT 1`

	u := &User{}

	err := m.DB.QueryRow(stmt).Scan(&u.ID, &u.ChatID, &u.IsSending, &u.RecipientID, pq.Array(&u.Blocks), &u.IsBan, &u.IsAdmin)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Ban changes the is_ban field of a user to true
func (m *UserModel) Ban(chatID int64) error {
	err := m.changeBanState(chatID, true)
	return err
}

// Unban changes the is_ban field of a user to false
func (m *UserModel) Unban(chatID int64) error {
	err := m.changeBanState(chatID, false)
	return err
}

// changeBanState is a local helper function that changes the is_ban field to true/false and used by Ban and Unban functions
func (m *UserModel) changeBanState(chatID int64, isBanning bool) error {
	stmt := `UPDATE users SET is_ban = $1
	WHERE chat_id = $2`

	_, err := m.DB.Exec(stmt, isBanning, chatID)
	if err != nil {
		return err
	}

	return nil
}

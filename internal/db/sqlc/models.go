// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Conversation struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type ConversationUser struct {
	ID             int64              `json:"id"`
	ConversationID int64              `json:"conversation_id"`
	UserID         int64              `json:"user_id"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
}

type Message struct {
	ID             int64              `json:"id"`
	UserID         int64              `json:"user_id"`
	Message        string             `json:"message"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	ConversationID int64              `json:"conversation_id"`
}

type User struct {
	ID                int64     `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	HashedPassword    string    `json:"hashed_password"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	IsEmailVerified   bool      `json:"is_email_verified"`
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateConversation(ctx context.Context, name string) (Conversation, error)
	CreateConversationUser(ctx context.Context, arg CreateConversationUserParams) (ConversationUser, error)
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteConversationUser(ctx context.Context, arg DeleteConversationUserParams) (ConversationUser, error)
	GetConversation(ctx context.Context, id int64) (Conversation, error)
	GetConversationMessages(ctx context.Context, arg GetConversationMessagesParams) ([]Message, error)
	GetConversationUser(ctx context.Context, id int64) (ConversationUser, error)
	GetConversationUsers(ctx context.Context, conversationID int64) ([]ConversationUser, error)
	GetMessage(ctx context.Context, id int64) (Message, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, email string) (User, error)
	GetUserConversations(ctx context.Context, arg GetUserConversationsParams) ([]Conversation, error)
	UpdateConversation(ctx context.Context, arg UpdateConversationParams) (Conversation, error)
	UpdateMessage(ctx context.Context, arg UpdateMessageParams) (Message, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)

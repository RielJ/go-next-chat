package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rielj/go-next-chat/internal/util"
)

func createRandomMessage(t *testing.T, conversationID int64, userID int64) Message {
	arg := CreateMessageParams{
		ConversationID: conversationID,
		UserID:         userID,
		Message:        util.RandomString(120),
	}

	message, err := testStore.CreateMessage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, message)

	require.Equal(t, arg.ConversationID, message.ConversationID)
	require.Equal(t, arg.UserID, message.UserID)
	require.Equal(t, arg.Message, message.Message)

	return message
}

func TestMessage(t *testing.T) {
	user := createRandomUser(t)
	c := createRandomConversation(t)
	m := createRandomMessage(t, c.ID, user.ID)

	require.NotEmpty(t, user)
	require.NotEmpty(t, c)
	require.NotEmpty(t, m)

	require.Equal(t, user.ID, m.UserID)
	require.Equal(t, c.ID, m.ConversationID)
}

func TestGroupMessages(t *testing.T) {
	c := createRandomConversation(t)
	user1 := createRandomUser(t)
	createRandomConversationUser(t, c.ID, user1.ID)
	m1 := createRandomMessage(t, c.ID, user1.ID)

	user2 := createRandomUser(t)
	createRandomConversationUser(t, c.ID, user2.ID)
	m2 := createRandomMessage(t, c.ID, user2.ID)

	params1 := GetConversationMessagesParams{
		ConversationID: c.ID,
		UserID:         user1.ID,
		Limit:          10,
		Offset:         0,
	}
	messages1, err := testStore.GetConversationMessages(context.Background(), params1)
	require.NoError(t, err)

	params2 := GetConversationMessagesParams{
		ConversationID: c.ID,
		UserID:         user2.ID,
		Limit:          10,
		Offset:         0,
	}
	messages2, err := testStore.GetConversationMessages(context.Background(), params2)
	require.NoError(t, err)

	require.NotEmpty(t, messages1)
	require.NotEmpty(t, messages2)

	require.Equal(t, m1.ID, messages1[0].ID)
	require.Equal(t, m2.ID, messages2[0].ID)
}

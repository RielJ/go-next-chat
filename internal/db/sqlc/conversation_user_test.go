package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomConversationUser(
	t *testing.T,
	conversationID int64,
	userId int64,
) ConversationUser {
	arg := CreateConversationUserParams{
		ConversationID: conversationID,
		UserID:         userId,
	}

	conversationUser, err := testStore.CreateConversationUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, conversationUser)

	require.Equal(t, arg.ConversationID, conversationUser.ConversationID)
	require.Equal(t, arg.UserID, conversationUser.UserID)

	return conversationUser
}

func TestConversationUser(t *testing.T) {
	user := createRandomUser(t)
	c := createRandomConversation(t)
	cu := createRandomConversationUser(t, c.ID, user.ID)

	require.NotEmpty(t, user)
	require.NotEmpty(t, c)
	require.NotEmpty(t, cu)

	require.Equal(t, user.ID, cu.UserID)
	require.Equal(t, c.ID, cu.ConversationID)
}

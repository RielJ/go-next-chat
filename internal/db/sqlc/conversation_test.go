package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rielj/go-next-chat/internal/util"
)

func createRandomConversation(t *testing.T) Conversation {
	randomName := util.RandomString(6)
	conversation, err := testStore.CreateConversation(context.Background(), randomName)
	require.NoError(t, err)
	require.NotEmpty(t, conversation)

	require.Equal(t, randomName, conversation.Name)

	return conversation
}

func TestConversation(t *testing.T) {
	user := createRandomUser(t)
	c := createRandomConversation(t)
	cu := createRandomConversationUser(t, c.ID, user.ID)
	m := createRandomMessage(t, c.ID, user.ID)

	require.NotEmpty(t, user)
	require.NotEmpty(t, c)
	require.NotEmpty(t, cu)
	require.NotEmpty(t, m)

	require.Equal(t, user.ID, cu.UserID)
	require.Equal(t, c.ID, cu.ConversationID)
	require.Equal(t, user.ID, m.UserID)
	require.Equal(t, c.ID, m.ConversationID)
}

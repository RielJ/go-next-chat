package db

import (
	"context"
	"fmt"
)

type MessageTxParams struct {
	FromUserID int64  `json:"from_user_id"`
	ToUserID   int64  `json:"to_user_id"`
	Message    string `json:"message"`
}

type MessageTxResult struct {
	FromUserID     int64  `json:"from_user_id"`
	ToUserID       int64  `json:"to_user_id"`
	Message        string `json:"message"`
	ConversationID int64  `json:"conversation_id"`
}

func (s *SQLStore) MessageTx(ctx context.Context, arg MessageTxParams) (MessageTxResult, error) {
	var result MessageTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		c, err := q.CreateConversation(ctx, fmt.Sprintf("%d:%d", arg.FromUserID, arg.ToUserID))
		_, err = q.CreateConversationUser(ctx, CreateConversationUserParams{
			UserID:         arg.FromUserID,
			ConversationID: c.ID,
		})
		_, err = q.CreateConversationUser(ctx, CreateConversationUserParams{
			UserID:         arg.ToUserID,
			ConversationID: c.ID,
		})

		_, err = q.CreateMessage(ctx, CreateMessageParams{
			ConversationID: c.ID,
			UserID:         arg.FromUserID,
			Message:        arg.Message,
		})
		if err != nil {
			return err
		}

		result = MessageTxResult{
			FromUserID:     arg.FromUserID,
			ToUserID:       arg.ToUserID,
			Message:        arg.Message,
			ConversationID: c.ID,
		}

		return nil
	})

	return result, err
}

package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	db "github.com/rielj/go-next-chat/internal/db/sqlc"
	"github.com/rielj/go-next-chat/internal/util"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

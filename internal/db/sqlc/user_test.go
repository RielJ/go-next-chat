package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	"github.com/rielj/go-next-chat/internal/util"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		HashedPassword: hashedPassword,
		FirstName:      util.RandomString(6),
		LastName:       util.RandomString(6),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
}

func TestUpdateUserOnlyName(t *testing.T) {
	oldUser := createRandomUser(t)

	arg := UpdateUserParams{
		Email: oldUser.Email,
		FirstName: pgtype.Text{
			String: util.RandomString(6),
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: util.RandomString(6),
			Valid:  true,
		},
	}

	user2, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, oldUser.ID, user2.ID)
	require.NotEqual(t, arg.FirstName, user2.FirstName)
	require.NotEqual(t, arg.LastName, user2.LastName)
	require.Equal(t, oldUser.HashedPassword, user2.HashedPassword)
	require.Equal(t, oldUser.Email, user2.Email)
	require.Equal(t, oldUser.PasswordChangedAt, user2.PasswordChangedAt)
	require.Equal(t, oldUser.CreatedAt, user2.CreatedAt)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	arg := UpdateUserParams{
		Email: oldUser.Email,
		HashedPassword: pgtype.Text{
			String: util.RandomString(6),
			Valid:  true,
		},
	}

	user2, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, oldUser.ID, user2.ID)
	require.Equal(t, oldUser.FirstName, user2.FirstName)
	require.Equal(t, oldUser.LastName, user2.LastName)
	require.NotEqual(t, oldUser.HashedPassword, user2.HashedPassword)
	require.Equal(t, oldUser.Email, user2.Email)
	require.Equal(t, oldUser.CreatedAt, user2.CreatedAt)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	arg := UpdateUserParams{
		Email: oldUser.Email,
		HashedPassword: pgtype.Text{
			String: util.RandomString(6),
			Valid:  true,
		},
		FirstName: pgtype.Text{
			String: util.RandomString(6),
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: util.RandomString(6),
			Valid:  true,
		},
		IsEmailVerified: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
	}

	user2, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, oldUser.ID, user2.ID)
	require.NotEqual(t, oldUser.FirstName, user2.FirstName)
	require.NotEqual(t, oldUser.LastName, user2.LastName)
	require.NotEqual(t, oldUser.HashedPassword, user2.HashedPassword)
	require.Equal(t, oldUser.Email, user2.Email)
	require.Equal(t, oldUser.CreatedAt, user2.CreatedAt)
	require.NotEqual(t, oldUser.IsEmailVerified, user2.IsEmailVerified)
}

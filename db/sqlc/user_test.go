package sqlc

import (
	"context"
	"database/sql"
	"github.com/borntodie-new/backend-master-class/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
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
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, newFullName, newUser.FullName)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)

}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, newEmail, newUser.Email)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)

}

func TestUpdateUserOnlyHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	password := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	ctime := time.Now()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  ctime,
			Valid: true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.WithinDuration(t, ctime, newUser.PasswordChangedAt, time.Second)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	newFullName := util.RandomOwner()
	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	ctime := time.Now()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  ctime,
			Valid: true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, newEmail, newUser.Email)
	require.NotEqual(t, oldUser.Email, newUser.Email)
	require.Equal(t, newFullName, newUser.FullName)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.WithinDuration(t, ctime, newUser.PasswordChangedAt, time.Second)
}

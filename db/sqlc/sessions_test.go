package sqlc

import (
	"context"
	"github.com/borntodie-new/backend-master-class/token"
	"github.com/borntodie-new/backend-master-class/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomSession(t *testing.T) (Session, uuid.UUID) {

	user := createRandomUser(t)

	jwtMaker, err := token.NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	username := util.RandomOwner()
	accessToken, accessPayload, err := jwtMaker.CreateToken(username, time.Minute)
	require.NoError(t, err)
	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: accessToken,
		UserAgent:    util.RandomString(12),
		ClientIp:     util.RandomString(12),
		IsBlock:      false,
		ExpiresAt:    accessPayload.ExpireAt,
	}
	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.Username, session.Username)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlock, session.IsBlock)

	require.NotZero(t, session.ExpiresAt)
	return session, arg.ID
}

func TestCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSession(t *testing.T) {
	session1, uid := createRandomSession(t)

	session2, err := testQueries.GetSession(context.Background(), uid)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.Username, session2.Username)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlock, session2.IsBlock)
	require.WithinDuration(t, session1.ExpiresAt, session2.ExpiresAt, time.Second)
}

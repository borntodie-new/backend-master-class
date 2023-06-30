package gapi

import (
	"context"
	"fmt"
	db "github.com/borntodie-new/backend-master-class/db/sqlc"
	"github.com/borntodie-new/backend-master-class/token"
	"github.com/borntodie-new/backend-master-class/util"
	"github.com/borntodie-new/backend-master-class/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	return server
}
func newContextBearerToken(t *testing.T, username string, duration time.Duration, tokenMaker token.Maker) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	bearer := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearer,
		},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}

package gapi

import (
	"context"
	"database/sql"
	db "github.com/borntodie-new/backend-master-class/db/sqlc"
	"github.com/borntodie-new/backend-master-class/pb"
	"github.com/borntodie-new/backend-master-class/util"
	"github.com/borntodie-new/backend-master-class/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if violations := validateLoginUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}
	// query this user
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %s", err)
	}

	// check the password
	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password: %s", err)
	}

	// sign the short token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %s", err)
	}

	// sign the long token refresh-token

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %s", err)
	}

	// extra metadata from ctx
	mtdt := s.extraMetadata(ctx)

	// create session
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlock:      false,
		ExpiresAt:    refreshPayload.ExpireAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpireAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpireAt),
		User:                  convertUser(user),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return
}

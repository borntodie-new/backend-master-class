package gapi

import (
	db "github.com/borntodie-new/backend-master-class/db/sqlc"
	"github.com/borntodie-new/backend-master-class/pb"
	"time"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: convertTime(user.PasswordChangedAt),
		CreatedAt:         convertTime(user.CreatedAt),
	}
}

func convertTime(t time.Time) *pb.Timestamp {
	return &pb.Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
}

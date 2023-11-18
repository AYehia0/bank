// maps the db objects to pb objects
package grpcapi

import (
	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUserToPb(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         timestamppb.New(user.CreatedAt),
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
	}
}

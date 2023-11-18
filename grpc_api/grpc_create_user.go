package grpcapi

import (
	"context"

	db "github.com/AYehia0/go-bk-mst/db/sqlc"
	"github.com/AYehia0/go-bk-mst/pb"
	"github.com/AYehia0/go-bk-mst/utils"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {

	// hashing the password
	hashedPassword, err := utils.GenerateHash(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "faild to hash password : %s", err)
	}

	arg := db.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		FullName: req.FullName,
	}
	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		// cast the pq error
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists : %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create a user : %s", err)
	}

	return convertUserToPb(user), nil
}

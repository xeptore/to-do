package user

import (
	"context"
	"database/sql"

	"github.com/xeptore/to-do/user/internal/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func New(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) VerifyPassword(ctx context.Context, in *pb.VerifyPasswordRequest) (*pb.VerifyPasswordReply, error) {
	// verify user email/password with database
	return nil, nil
}

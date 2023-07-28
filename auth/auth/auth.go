package auth

import (
	"context"
	"fmt"

	"github.com/xeptore/to-do/auth/internal/pb"
	"github.com/xeptore/to-do/auth/jwt"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	j jwt.JWT
	u pb.UserServiceClient
}

func New(j jwt.JWT, u pb.UserServiceClient) *AuthService {
	return &AuthService{j: j, u: u}
}

func (s *AuthService) VerifyToken(ctx context.Context, in *pb.VerifyTokenRequest) (*pb.VerifyTokenReply, error) {
	res, err := s.j.VerifyToken(ctx, in.Token)
	if nil != err {
		return nil, fmt.Errorf("auth: failed to verify token: %v", err)
	}
	return &pb.VerifyTokenReply{User: &pb.VerifyTokenReply_User{Id: res.UserID}}, nil
}

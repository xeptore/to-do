package auth

import (
	"context"
	"fmt"

	pbauth "github.com/xeptore/to-do/api/pb/auth"
	pbuser "github.com/xeptore/to-do/api/pb/user"
	"github.com/xeptore/to-do/auth/jwt"
)

type AuthService struct {
	pbauth.UnimplementedAuthServiceServer
	j jwt.JWT
	u pbuser.UserServiceClient
}

func New(j jwt.JWT, u pbuser.UserServiceClient) *AuthService {
	return &AuthService{j: j, u: u}
}

func (s *AuthService) VerifyToken(ctx context.Context, in *pbauth.VerifyTokenRequest) (*pbauth.VerifyTokenReply, error) {
	res, err := s.j.VerifyToken(ctx, in.Token)
	if nil != err {
		return nil, fmt.Errorf("auth: failed to verify token: %v", err)
	}
	return &pbauth.VerifyTokenReply{User: &pbauth.VerifyTokenReply_User{Id: res.UserID}}, nil
}

package auth

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"

	"github.com/xeptore/to-do/api/pb"
)

func (s *AuthService) HandleMessage(ctx context.Context, msg *nats.Msg) []byte {
	command := gjson.GetBytes(msg.Data, "command").String()
	switch command {
	case "login":
		return s.Login(ctx, gjson.GetBytes(msg.Data, "payload").Raw)
	default:
		return []byte("unknown command")
	}
}

type LoginRequest struct {
	Email    string
	Name     string
	Password string
}

type LoginResult struct {
	Token string
}

func (s *AuthService) Login(ctx context.Context, in string) []byte {
	var req LoginRequest
	if err := json.UnmarshalContext(ctx, []byte(in), &req); nil != err {
		return []byte("invalid request")

	}

	res, err := s.u.VerifyPassword(ctx, &pb.VerifyPasswordRequest{Email: req.Email, Password: req.Password})
	if nil != err {
		return []byte("user email/password verification failed")
	}

	token, err := s.j.GenerateToken(ctx, res.User.Id)
	if nil != err {
		return []byte("failed to generate jwt token")
	}

	out, err := json.MarshalContext(ctx, LoginResult{Token: token.Token})
	if nil != err {
		return []byte("failed to prepare response")
	}

	return out
}

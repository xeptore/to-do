package user

import (
	"context"

	"github.com/goccy/go-json"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func (s *UserService) HandleMessage(ctx context.Context, msg *nats.Msg) []byte {
	command := gjson.GetBytes(msg.Data, "command").String()
	switch command {
	case "create_user":
		return s.Create(ctx, gjson.GetBytes(msg.Data, "payload").Raw)
	default:
		return []byte("unknown command")
	}
}

type CreateRequest struct {
	Email    string
	Name     string
	Password string
}

type CreateResult struct {
	UserID string
}

func (s *UserService) Create(ctx context.Context, in string) []byte {
	var req CreateRequest
	if err := json.UnmarshalContext(ctx, []byte(in), &req); nil != err {
		return []byte("invalid request")
	}
	userID, err := gonanoid.New(24)
	if nil != err {
		return []byte("failed to generate userID")
	}
	// insert into database

	res := CreateResult{UserID: userID}
	out, err := json.MarshalContext(ctx, res)
	if nil != err {
		return []byte("failed to prepare response")
	}

	return out
}

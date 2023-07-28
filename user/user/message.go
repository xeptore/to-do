package user

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"

	m "github.com/xeptore/to-do/user/db/gen/todo/public/model"
	t "github.com/xeptore/to-do/user/db/gen/todo/public/table"
)

func (s *UserService) HandleMessage(ctx context.Context, msg *nats.Msg) []byte {
	command := gjson.GetBytes(msg.Data, "command").String()
	switch command {
	// TODO: export available commands as constants that can be imported by other services.
	case "create_user":
		return s.Create(ctx, gjson.GetBytes(msg.Data, "payload").Raw)
	default:
		// TODO: return valid typed error that can be handled by the client
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
		// TODO: return typed validation error that can be handled by the client
		return []byte("invalid request")
	}

	userID, err := gonanoid.New(16)
	if nil != err {
		// TODO: return typed internal error that can be handled by the client
		// TODO: log the error
		return []byte("failed to generate userID")
	}

	// FIXME: store hashed passwords
	model := m.Users{
		ID:          userID,
		TheName:     req.Name,
		Email:       req.Email,
		ThePassword: req.Password,
		CreatedAt:   time.Now(),
	}
	if res, err := t.Users.INSERT(t.Users.AllColumns).MODEL(model).ExecContext(ctx, s.db); nil != err {
		// TODO: handle duplicate email error, and return validation error
		// TODO: return typed internal error that can be handled by the client
		// TODO: log the error
		return []byte("failed to execute user insert statement")
	} else if affectedRows, err := res.RowsAffected(); nil != err {
		// TODO: return typed internal error that can be handled by the client
		// TODO: log the error
		return []byte("failed to get number of affected rows by user insert query execution")
	} else if affectedRows != 1 {
		// TODO: return typed internal error that can be handled by the client
		// TODO: log the error
		return []byte("expected only 1 row to be affected by user insert query execution")
	}

	res := CreateResult{UserID: userID}
	out, err := json.MarshalContext(ctx, res)
	if nil != err {
		// TODO: return typed internal error that can be handled by the client
		// TODO: log the error
		return []byte("failed to prepare response")
	}

	return out
}

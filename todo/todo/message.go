package todo

import (
	"context"

	"github.com/nats-io/nats.go"
)

func (s *TodoService) HandleMessage(ctx context.Context, msg *nats.Msg) []byte {
	// TODO
	return nil
}

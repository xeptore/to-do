package todo

import (
	"context"

	"github.com/nats-io/nats.go"
)

func (s *TodoService) HandleMessage(ctx context.Context, msg *nats.Msg) []byte {
	// TODO: implement the nats on-message handler that distributes the command to its respective handler.
	return nil
}

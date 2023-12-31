package todo

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	pb "github.com/xeptore/to-do/api/pb/todo"

	m "github.com/xeptore/to-do/todo/db/gen/sample/public/model"
	t "github.com/xeptore/to-do/todo/db/gen/sample/public/table"
)

type TodoService struct {
	pb.UnimplementedTodoServiceServer
	db *sql.DB
}

func New(db *sql.DB) *TodoService {
	return &TodoService{db: db}
}

func (s *TodoService) GetList(ctx context.Context, in *pb.GetListRequest) (*pb.GetListReply, error) {
	var model m.TodoLists
	// FIXME: this needs to be moved to its own repository package
	err := t.TodoLists.SELECT(t.TodoLists.AllColumns).WHERE(t.TodoLists.ID.EQ(postgres.String(in.Id))).QueryContext(ctx, s.db, &model)
	if nil != err {
		// TODO: handle not found error
		// TODO: log internal error
		// TODO: reply with grpc-compatible internal error response
	}

	return &pb.GetListReply{
		Id:          model.ID,
		Name:        model.TheName,
		Description: model.TheDescription,
		CreatedById: model.CreatedByID,
	}, nil
}

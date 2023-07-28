package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbuser "github.com/xeptore/to-do/api/pb/user"
	m "github.com/xeptore/to-do/user/db/gen/todo/public/model"
	t "github.com/xeptore/to-do/user/db/gen/todo/public/table"
)

type UserService struct {
	pbuser.UnimplementedUserServiceServer
	db *sql.DB
}

func New(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) VerifyPassword(ctx context.Context, in *pbuser.VerifyPasswordRequest) (*pbuser.VerifyPasswordReply, error) {
	var model m.Users
	stmt := t.Users.
		SELECT(t.Users.ID).
		WHERE(
			postgres.AND(
				t.Users.Email.EQ(postgres.String(in.Email)),
				t.Users.ThePassword.EQ(postgres.String(in.Password)),
			),
		)
	if err := stmt.QueryContext(ctx, s.db, &model); nil != err {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user was not found")
		}
		return nil, status.Error(codes.Internal, "failed to query user")
	}
	return &pbuser.VerifyPasswordReply{User: &pbuser.VerifyPasswordReply_User{Id: model.ID}}, nil
}

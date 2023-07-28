package api

import (
	"net/http"

	"github.com/goccy/go-json"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/nats.go"
	pbauth "github.com/xeptore/to-do/api/pb/auth"
	pbtodo "github.com/xeptore/to-do/api/pb/todo"
	pbuser "github.com/xeptore/to-do/api/pb/user"
	"github.com/xeptore/to-do/auth/auth"
	"github.com/xeptore/to-do/user/user"
)

type Server struct {
	userGrpc pbuser.UserServiceClient
	authGrpc pbauth.AuthServiceClient
	todoGrpc pbtodo.TodoServiceClient
	userNats *nats.EncodedConn
	authNats *nats.EncodedConn
}

func NewServer(
	userGrpc pbuser.UserServiceClient,
	authGrpc pbauth.AuthServiceClient,
	todoGrpc pbtodo.TodoServiceClient,
	userNats *nats.EncodedConn,
	authNats *nats.EncodedConn,
) *Server {
	return &Server{userGrpc: userGrpc, authGrpc: authGrpc, todoGrpc: todoGrpc, userNats: userNats, authNats: authNats}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req user.CreateRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &req); nil != err {
		// TODO: respond with validation error
		return
	}

	// FIXME: "command", and "payload" parameters must be sent as the request data
	var res user.CreateResult
	if err := s.userNats.RequestWithContext(ctx, "request", req, &res); nil != err {
		// TODO: respond with internal error
		// TODO: log the error
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalContext(ctx, res)
	if nil != err {
		// TODO: log the error
	}
	if n, err := w.Write(b); nil != err {
		// TODO: log the error
	} else if n != len(b) {
		// TODO: log the error
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req auth.LoginRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &req); nil != err {
		// TODO: respond with validation error
		return
	}

	// FIXME: "command", and "payload" parameters must be sent as the request data
	var res auth.LoginResult
	if err := s.authNats.RequestWithContext(ctx, "request", req, &res); nil != err {
		// TODO: respond with internal error
		// TODO: log the error
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalContext(ctx, res)
	if nil != err {
		// TODO: log the error
	}
	if n, err := w.Write(b); nil != err {
		// TODO: log the error
	} else if n != len(b) {
		// TODO: log the error
	}
}

type GetTodoListResult struct {
	ID          string
	Name        string
	Description string
	CreatedByID string
}

func (s *Server) GetTodoList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	// TODO: add validation

	ctx := r.Context()
	res, err := s.todoGrpc.GetList(ctx, &pbtodo.GetListRequest{Id: id})
	if nil != err {
		// TODO: respond with internal error
		// TODO: log the error
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	out := GetTodoListResult{
		ID:          res.Id,
		Name:        res.Name,
		Description: res.Description,
		CreatedByID: res.CreatedById,
	}
	b, err := json.MarshalContext(ctx, out)
	if nil != err {
		// TODO: log the error
	}
	if n, err := w.Write(b); nil != err {
		// TODO: log the error
	} else if n != len(b) {
		// TODO: log the error
	}
}

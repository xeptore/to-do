package api

import (
	"net/http"

	"github.com/goccy/go-json"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/nats.go"
	"github.com/xeptore/to-do/auth/auth"
	"github.com/xeptore/to-do/user/user"

	"github.com/xeptore/to-do/gateway/internal/pb"
)

type Server struct {
	userGrpc pb.UserServiceClient
	authGrpc pb.AuthServiceClient
	userNats *nats.EncodedConn
	authNats *nats.EncodedConn
}

func NewServer(
	userGrpc pb.UserServiceClient,
	authGrpc pb.AuthServiceClient,
	userNats *nats.EncodedConn,
	authNats *nats.EncodedConn,
) *Server {
	return &Server{userGrpc: userGrpc, authGrpc: authGrpc, userNats: userNats, authNats: authNats}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req user.CreateRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &req); nil != err {
		// TODO: handle error
	}
	var res user.CreateResult
	if err := s.userNats.RequestWithContext(ctx, "request", req, &res); nil != err {
		// TODO: handle error
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalContext(ctx, res)
	if nil != err {
		// TODO: handle error
	}
	if n, err := w.Write(b); nil != err {
		// TODO: handle error
	} else if n != len(b) {
		// TODO: handle error
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req auth.LoginRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).DecodeContext(ctx, &req); nil != err {
		// TODO: handle error
	}

	var res auth.LoginResult
	if err := s.authNats.RequestWithContext(ctx, "request", req, &res); nil != err {
		// TODO: handle error
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := json.MarshalContext(ctx, res)
	if nil != err {
		// TODO: handle error
	}
	if n, err := w.Write(b); nil != err {
		// TODO: handle error
	} else if n != len(b) {
		// TODO: handle error
	}
}

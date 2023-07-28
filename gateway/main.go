package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/xeptore/to-do/config"
	"google.golang.org/grpc"

	"github.com/xeptore/to-do/gateway/api"
	"github.com/xeptore/to-do/gateway/internal/pb"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) { w.Out = os.Stderr; w.TimeFormat = time.RFC3339 })).With().Timestamp().Logger().Level(zerolog.TraceLevel)

	f, err := os.Open("config.yml")
	if nil != err {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal().Err(err).Msg("failed to load config.yml")
		}
		log.Warn().Msg("config file was not found")
	}
	cfg, err := config.FromYaml(ctx, f)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to load configuration from config file")
	}

	if err := godotenv.Load(".env"); nil != err {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal().Err(err).Msg("unexpected error while loading environment variables from .env file")
		}
		log.Warn().Msg(".env file not found")
	}

	tz, ok := os.LookupEnv("TZ")
	if !ok || tz != "UTC" {
		log.Fatal().Msg("TZ environment variable must be set to UTC")
	}

	userNatsConn, err := nats.Connect(cfg.Str("user.nats_address"))
	if nil != err {
		log.Fatal().Err(err).Msg("failed to connect to user nats server")
	}
	userNatsClient, err := nats.NewEncodedConn(userNatsConn, nats.JSON_ENCODER)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to initialize user nats json-encoded connection")
	}

	authNatsConn, err := nats.Connect(cfg.Str("auth.nats_address"))
	if nil != err {
		log.Fatal().Err(err).Msg("failed to connect to auth nats server")
	}
	authNatsClient, err := nats.NewEncodedConn(authNatsConn, nats.JSON_ENCODER)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to initialize auth nats json-encoded connection")
	}

	userGrpcConn, err := grpc.Dial(cfg.Str("user.grpc_address"), []grpc.DialOption{grpc.WithInsecure()}...)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to connect to user grpc service")
	}
	userGrpcClient := pb.NewUserServiceClient(userGrpcConn)

	authGrpcConn, err := grpc.Dial(cfg.Str("auth.grpc_address"), []grpc.DialOption{grpc.WithInsecure()}...)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to connect to auth grpc service")
	}
	authGrpcClient := pb.NewAuthServiceClient(authGrpcConn)

	srv := api.NewServer(userGrpcClient, authGrpcClient, userNatsClient, authNatsClient)

	router := httprouter.New()
	router.POST("/users", srv.CreateUser)
	router.POST("/auth/login", srv.Login)
	httpServer := http.Server{Addr: cfg.Str("server.listen_addr"), Handler: router}

	done := make(chan bool)
	go func() {
		<-ctx.Done()
		log.Info().Msg("executing cleanup tasks as on root context cancellation...")
		if err := httpServer.Shutdown(ctx); nil != err {
			log.Error().Err(err).Msg("failed to gracefully shutdown http server")
		}
		userNatsConn.Close()
		userNatsClient.Close()
		authNatsConn.Close()
		authNatsClient.Close()
		if err := authGrpcConn.Close(); nil != err {
			log.Error().Err(err).Msg("failed to gracefully close auth service grpc client connection")
		}
		if err := userGrpcConn.Close(); nil != err {
			log.Error().Err(err).Msg("failed to gracefully close user service grpc client connection")
		}
	}()

	if err := httpServer.ListenAndServe(); nil != err {
		cancel()
		<-done
		log.Fatal().Err(err).Msg("failed to start http server")
	}
	<-done
}

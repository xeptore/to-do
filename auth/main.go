package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	pbauth "github.com/xeptore/to-do/api/pb/auth"
	pbuser "github.com/xeptore/to-do/api/pb/user"
	"github.com/xeptore/to-do/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/xeptore/to-do/auth/auth"
	"github.com/xeptore/to-do/auth/jwt"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) { w.Out = os.Stderr; w.TimeFormat = time.RFC3339 })).With().Timestamp().Logger().Level(zerolog.TraceLevel)

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

	userGrpcDialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	userGrpcConn, err := grpc.Dial(cfg.Str("user.grpc_address"), userGrpcDialOpts...)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to connect to user grpc service")
	}
	userGrpcClient := pbuser.NewUserServiceClient(userGrpcConn)

	jeydubti := jwt.New([]byte(os.Getenv("JWT_SECRET")))
	authService := auth.New(jeydubti, userGrpcClient)
	grpcSrv := grpc.NewServer(grpc.ConnectionTimeout(time.Second*3), grpc.MaxConcurrentStreams(10))
	pbauth.RegisterAuthServiceServer(grpcSrv, authService)

	lis, err := net.Listen("tcp", cfg.Str("grpc.listen_address"))
	if nil != err {
		log.Fatal().Err(err).Msg("failed to bind grpc server to address")
	}

	conn, err := nats.Connect(cfg.Str("nats.address"))
	if nil != err {
		log.Fatal().Err(err).Msg("failed to connect to nats server")
	}
	nc, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to initialize nats json-encoded connection")
	}

	sub, err := nc.Subscribe("request", func(msg *nats.Msg) {
		resp := authService.HandleMessage(ctx, msg)
		if err := msg.Respond(resp); nil != err {
			log.Error().Err(err).Bytes("response_message_data", resp).Msg("failed to respond to message")
		}
	})
	if nil != err {
		log.Fatal().Err(err).Msg("failed to subscribe to nats request stream")
	}

	done := make(chan bool)
	go func() {
		<-ctx.Done()
		log.Info().Msg("executing cleanup tasks as on root context cancellation...")
		if err := sub.Unsubscribe(); nil != err {
			log.Error().Err(err).Msg("failed to unsubscribe nats subscription")
		}
		if err := sub.Drain(); nil != err {
			log.Error().Err(err).Msg("failed to drain nats subscription")
		}
		if err := userGrpcConn.Close(); nil != err {
			log.Error().Err(err).Msg("failed to close user service grpc client connection")
		}
		grpcSrv.GracefulStop()
		done <- true
	}()

	if err := grpcSrv.Serve(lis); nil != err {
		cancel()
		<-done
		log.Fatal().Err(err).Msg("failed to start grpc server")
	}
	<-done
}

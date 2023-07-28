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
	"google.golang.org/grpc"

	"github.com/xeptore/to-do/config"

	"github.com/xeptore/to-do/user/db"
	"github.com/xeptore/to-do/user/internal/pb"
	"github.com/xeptore/to-do/user/user"
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

	database, err := db.Connect(ctx)
	if nil != err {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	userService := user.New(database)
	grpcSrv := grpc.NewServer(grpc.ConnectionTimeout(time.Second*3), grpc.MaxConcurrentStreams(10))
	pb.RegisterUserServiceServer(grpcSrv, userService)
	lis, err := net.Listen("tcp", cfg.Str("grpc.listen_address"))
	if nil != err {
		log.Fatal().Err(err).Msg("failed to bind grpc server to address")
	}

	nc, err := nats.Connect(cfg.Str("nats.address"))
	if nil != err {
		log.Fatal().Err(err).Msg("failed to connect to nats server")
	}

	sub, err := nc.Subscribe("request", func(msg *nats.Msg) {
		resp := userService.HandleMessage(ctx, msg)
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

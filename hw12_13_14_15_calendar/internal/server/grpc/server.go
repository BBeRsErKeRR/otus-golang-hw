package internalgrpc

import (
	"context"
	"log"
	"net"

	router "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api"
	event "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api/v1/grpc"
	v1grpc "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api/v1/grpc"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	logger logger.Logger
	Addr   string
	server *grpc.Server
}

type Config struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func NewServer(logger logger.Logger, app router.Application, conf *Config) *Server {
	addr, err := server.GetAddress(conf.Host, conf.Port)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			loggingMiddleware(logger),
		),
	)
	event.RegisterEventServiceServer(server, v1grpc.NewHandler(app, logger))
	return &Server{
		Addr:   addr,
		logger: logger,
		server: server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	list, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.logger.Info("starting server", zap.String("address", s.Addr))
	return s.server.Serve(list)
}

func (s *Server) Stop() error {
	s.server.GracefulStop()
	return nil
}

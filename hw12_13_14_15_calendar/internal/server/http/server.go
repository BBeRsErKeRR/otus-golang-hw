package internalhttp

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	httprouter "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api"
	v1routes "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api/v1/http"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/server"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ErrorUnsupportedHostAddress = errors.New("invalid host address")
	ErrorInvalidPort            = errors.New("invalid port number")
)

type Server struct {
	server *http.Server
	logger logger.Logger
}

type Config struct {
	Host              string        `mapstructure:"host"`
	Port              string        `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
}

func NewServer(logger logger.Logger, app httprouter.Application, conf *Config) *Server {
	addr, err := server.GetAddress(conf.Host, conf.Port)
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	handlerV1 := v1routes.NewHandler(app, logger)
	handlerV1.AddV1Routes(router)
	return &Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           loggingMiddleware(router, logger),
			ReadTimeout:       conf.ReadTimeout,
			WriteTimeout:      conf.WriteTimeout,
			ReadHeaderTimeout: conf.ReadHeaderTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("starting server", zap.String("address", s.server.Addr))
	err := s.server.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return s.server.Shutdown(ctx)
}

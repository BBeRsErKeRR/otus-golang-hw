package internalhttp

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
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

type Application interface {
}

func getAddress(hostArg string, portArg string) (string, error) {
	var address string
	re := regexp.MustCompile(`^((([a-z0-9][a-z0-9\-]*[a-z0-9])|[a-z0-9])\.?)+$`)

	if hostArg != "localhost" && !re.MatchString(hostArg) && net.ParseIP(hostArg) == nil {
		return address, ErrorUnsupportedHostAddress
	}

	port, err := strconv.Atoi(portArg)
	if err != nil {
		return address, err
	}

	if port < 1 || port > 65535 {
		return address, ErrorInvalidPort
	}

	address = net.JoinHostPort(hostArg, portArg)
	return address, nil
}

func NewServer(logger logger.Logger, app Application, conf *Config) *Server {
	addr, err := getAddress(conf.Host, conf.Port)
	if err != nil {
		log.Fatal(err)
	}
	return &Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           loggingMiddleware(NewHandler(app, logger), logger),
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

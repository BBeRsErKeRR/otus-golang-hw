package internalhttp

import (
	"net"
	"net/http"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"go.uber.org/zap"
)

type CatchStatusResponseWriter struct {
	http.ResponseWriter
	Status int
}

func loggingMiddleware(next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		csr := &CatchStatusResponseWriter{ResponseWriter: w}
		next.ServeHTTP(csr, r)
		logger.Info(
			"access-log",
			zap.String("ip", ip),
			// zap.String("datetime", time.Now().Format(time.RFC822)),
			zap.String("method", r.Method),
			zap.String("path", r.URL.EscapedPath()),
			zap.String("protocol_version", r.Proto),
			zap.String("user-agent", r.UserAgent()),
			zap.Any("duration", time.Since(start).String()),
			zap.Any("http-status", csr.Status),
		)
	})
}

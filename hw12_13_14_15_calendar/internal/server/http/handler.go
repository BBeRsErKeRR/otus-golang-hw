package internalhttp

import (
	"net/http"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"go.uber.org/zap"
)

type Handler struct {
	app    Application
	logger logger.Logger
}

func NewHandler(app Application, logger logger.Logger) *Handler {
	return &Handler{
		app:    app,
		logger: logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.String() {
	case "/hello":
		h.helloWorld(w, r)
	default:
		h.notFound(w, r)
	}
}

func (h *Handler) sendResponse(data string, status int, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write([]byte(data)); err != nil {
		status = http.StatusInternalServerError
		h.logger.Error("error send data to client", zap.Error(err))
	}
	w.WriteHeader(status)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) { //nolint:unparam
	h.sendResponse(`{"error":"Not found"}"`, http.StatusNotFound, w)
}

func (h *Handler) helloWorld(w http.ResponseWriter, r *http.Request) { //nolint:unparam
	h.sendResponse(`{"msg":"Hello, world!"}"`, http.StatusOK, w)
}

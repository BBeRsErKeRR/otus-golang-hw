package v1routes

import (
	"net/http"

	httprouter "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	app    httprouter.Application
	logger logger.Logger
}

func NewHandler(app httprouter.Application, logger logger.Logger) *Handler {
	return &Handler{
		app:    app,
		logger: logger,
	}
}

func (h *Handler) AddV1Routes(r *mux.Router) {
	h.addRoutes(r.PathPrefix("/v1").Subrouter())
	r.NotFoundHandler = http.HandlerFunc(h.notFound)
}

func (h *Handler) addRoutes(r *mux.Router) {
	r.HandleFunc("/hello", h.helloWorld)
}

func (h *Handler) sendResponse(data string, status int, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write([]byte(data)); err != nil {
		status = http.StatusInternalServerError
		h.logger.Error("error send data to client", zap.Error(err))
	}
	w.WriteHeader(status)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"error":"Not found"}"`, http.StatusNotFound, w)
}

func (h *Handler) helloWorld(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"msg":"Hello, world!"}"`, http.StatusOK, w)
}

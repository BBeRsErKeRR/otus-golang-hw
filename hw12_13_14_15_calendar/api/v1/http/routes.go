package v1routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	router "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type EventDTO struct {
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	EndDate    time.Time `json:"end_date"`
	Desc       string    `json:"description"`
	UserID     string    `json:"user_id"`
	RemindDate time.Time `json:"remind_date"`
}

func (e *EventDTO) transfer() storage.Event {
	return storage.Event{
		Title:      e.Title,
		Date:       e.Date,
		EndDate:    e.EndDate,
		Desc:       e.Desc,
		UserID:     e.UserID,
		RemindDate: e.RemindDate,
	}
}

type Handler struct {
	app    router.Application
	logger logger.Logger
}

func NewHandler(app router.Application, logger logger.Logger) *Handler {
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
	r.HandleFunc("/create", h.CreateEvent)
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

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.sendResponse(fmt.Sprintf(`{"error":"%v"}"`, err), http.StatusBadRequest, w)
		return
	}
	var data EventDTO
	json.Unmarshal(reqBody, &data)
	err = h.app.CreateEvent(r.Context(), data.transfer())
	if err != nil {
		h.sendResponse(fmt.Sprintf(`{"error":"%v"}"`, err), http.StatusBadRequest, w)
		return
	}
	h.sendResponse(`{"msg":"Created"}`, http.StatusOK, w)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"msg":"Hello, world!"}"`, http.StatusOK, w)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"msg":"Hello, world!"}"`, http.StatusOK, w)
}

func (h *Handler) GetDailyEvents(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"msg":"Hello, world!"}"`, http.StatusOK, w)
}

func (h *Handler) GetWeeklyEvents(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"msg":"Hello, world!"}"`, http.StatusOK, w)
}

func (h *Handler) GetMonthlyEvents(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"msg":"Hello, world!"}"`, http.StatusOK, w)
}

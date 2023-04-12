package v1routes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	router "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var ErrBadArgs = errors.New("bad args")

type DateRequest struct {
	Date time.Time `json:"date"`
}

type EventDTO struct {
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	EndDate    time.Time `json:"end_date"` //nolint:tagliatelle
	Desc       string    `json:"description"`
	UserID     string    `json:"user_id"`     //nolint:tagliatelle
	RemindDate time.Time `json:"remind_date"` //nolint:tagliatelle
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
	r.HandleFunc("/hello", h.helloWorld).Methods("GET")
	r.HandleFunc("/event", h.CreateEvent).Methods("POST")
	r.HandleFunc("/event/{id}", h.UpdateEvent).Methods("PUT")
	r.HandleFunc("/event/{id}", h.DeleteEvent).Methods("DELETE")
	r.HandleFunc("/events/daily", h.GetDailyEvents).Methods("GET", "POST")
	r.HandleFunc("/events/weekly", h.GetWeeklyEvents).Methods("GET", "POST")
	r.HandleFunc("/events/monthly", h.GetMonthlyEvents).Methods("GET", "POST")
}

func (h *Handler) sendResponse(data []byte, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		status = http.StatusInternalServerError
		w.WriteHeader(status)
		h.logger.Error("error send data to client", zap.Error(err))
	}
}

func (h *Handler) sendError(err error, status int, w http.ResponseWriter) { //nolint:unparam
	h.sendResponse([]byte(fmt.Sprintf(`{"error":"%v"}`, err)), status, w)
}

func (h *Handler) sendString(data string, status int, w http.ResponseWriter) {
	h.sendResponse([]byte(data), status, w)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.sendString(`{"error":"Not found"}"`, http.StatusNotFound, w)
}

func (h *Handler) helloWorld(w http.ResponseWriter, r *http.Request) {
	h.sendString(`{"msg":"Hello, world!"}"`, http.StatusOK, w)
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	var data EventDTO
	json.Unmarshal(reqBody, &data)
	id, err := h.app.CreateEvent(r.Context(), data.transfer())
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	h.sendString(fmt.Sprintf(`{"msg":"Created","id":"%v"}`, id), http.StatusCreated, w)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		h.sendError(ErrBadArgs, http.StatusBadRequest, w)
		return
	}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	var data EventDTO
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	err = h.app.UpdateEvent(r.Context(), id, data.transfer())
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	h.sendString(`{"msg":"Updated"}`, http.StatusOK, w)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		h.sendError(ErrBadArgs, http.StatusBadRequest, w)
		return
	}
	err := h.app.DeleteEvent(r.Context(), id)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	h.sendString(`{"msg":"Deleted"}`, http.StatusOK, w)
}

func (h *Handler) getDateParams(r *http.Request) (time.Time, error) {
	var res time.Time
	dateStr := r.URL.Query().Get("date")
	if dateStr != "" {
		return time.Parse(time.RFC3339, dateStr)
	}
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return res, err
	}
	var data DateRequest
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		return res, err
	}

	return data.Date, nil
}

type EventsResponse struct {
	Events []storage.Event `json:"events"`
}

func (h *Handler) GetDailyEvents(w http.ResponseWriter, r *http.Request) {
	date, err := h.getDateParams(r)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	events, err := h.app.GetDailyEvents(r.Context(), date)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	res, err := json.Marshal(EventsResponse{
		Events: events,
	})
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	h.sendResponse(res, http.StatusOK, w)
}

func (h *Handler) GetWeeklyEvents(w http.ResponseWriter, r *http.Request) {
	date, err := h.getDateParams(r)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	events, err := h.app.GetWeeklyEvents(r.Context(), date)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	res, err := json.Marshal(EventsResponse{
		Events: events,
	})
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	h.sendResponse(res, http.StatusOK, w)
}

func (h *Handler) GetMonthlyEvents(w http.ResponseWriter, r *http.Request) {
	date, err := h.getDateParams(r)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	events, err := h.app.GetMonthlyEvents(r.Context(), date)
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	res, err := json.Marshal(EventsResponse{
		Events: events,
	})
	if err != nil {
		h.sendError(err, http.StatusBadRequest, w)
		return
	}
	h.sendResponse(res, http.StatusOK, w)
}

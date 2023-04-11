package v1grpc

import (
	"context"
	"time"

	router "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	app    router.Application
	logger logger.Logger
	UnimplementedEventServiceServer
}

func NewHandler(app router.Application, logger logger.Logger) *Handler {
	return &Handler{
		app:    app,
		logger: logger,
	}
}

func (h *Handler) getStorageEvent(event *Event) storage.Event {
	newEvent := storage.Event{
		Title:  event.GetTitle(),
		Desc:   event.GetDesc(),
		UserID: event.GetUserID(),
	}
	if event.GetDate() != nil {
		newEvent.Date = event.GetDate().AsTime()
	}
	if event.GetEndDate() != nil {
		newEvent.EndDate = event.GetEndDate().AsTime()
	}
	if event.GetRemindDate() != nil {
		newEvent.RemindDate = event.GetRemindDate().AsTime()
	}
	return newEvent
}

func (h *Handler) getRequestEvents(events []storage.Event) []*Event {
	res := make([]*Event, len(events))
	for i, el := range events {
		event := Event{
			ID:         el.ID,
			Title:      el.Title,
			Date:       timestamppb.New(el.Date),
			EndDate:    timestamppb.New(el.EndDate),
			Desc:       el.Desc,
			UserID:     el.UserID,
			RemindDate: timestamppb.New(el.RemindDate),
		}
		res[i] = &event
	}
	return res
}

func (h *Handler) CreateEvent(ctx context.Context, event *Event) (*EventIDResponse, error) {
	id, err := h.app.CreateEvent(ctx, h.getStorageEvent(event))
	if err != nil {
		return &EventIDResponse{Msg: "error"}, err
	}
	return &EventIDResponse{Msg: "Created", Id: id}, nil
}

func (h *Handler) UpdateEvent(ctx context.Context, req *UpdateRequest) (*Response, error) {
	err := h.app.UpdateEvent(ctx, req.Id, h.getStorageEvent(req.Event))
	if err != nil {
		return &Response{Msg: "error"}, err
	}
	return &Response{Msg: "success"}, nil
}

func (h *Handler) DeleteEvent(ctx context.Context, req *EventID) (*Response, error) {
	h.logger.Info("DeleteEvent", zap.String("eventID", req.Id))
	err := h.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		return &Response{Msg: "error"}, err
	}
	return &Response{Msg: "success"}, nil
}

func (h *Handler) GetDailyEvents(ctx context.Context, req *EventsRequest) (*EventsResponse, error) {
	date := req.Date.AsTime()
	h.logger.Debug("GetDailyEvents", zap.String("datetime", date.Format(time.RFC822)))
	events, err := h.app.GetDailyEvents(ctx, date)
	if err != nil {
		return &EventsResponse{Error: err.Error()}, err
	}
	return &EventsResponse{Events: h.getRequestEvents(events)}, nil
}

func (h *Handler) GetWeeklyEvents(ctx context.Context, req *EventsRequest) (*EventsResponse, error) {
	date := req.Date.AsTime()
	h.logger.Debug("GetWeeklyEvents", zap.String("datetime", date.Format(time.RFC822)))
	events, err := h.app.GetWeeklyEvents(ctx, date)
	if err != nil {
		return &EventsResponse{Error: err.Error()}, err
	}
	return &EventsResponse{Events: h.getRequestEvents(events)}, nil
}

func (h *Handler) GetMonthlyEvents(ctx context.Context, req *EventsRequest) (*EventsResponse, error) {
	date := req.Date.AsTime()
	h.logger.Debug("GetMonthlyEvents", zap.String("datetime", date.Format(time.RFC822)))
	events, err := h.app.GetMonthlyEvents(ctx, date)
	if err != nil {
		return &EventsResponse{Error: err.Error()}, err
	}
	return &EventsResponse{Events: h.getRequestEvents(events)}, nil
}

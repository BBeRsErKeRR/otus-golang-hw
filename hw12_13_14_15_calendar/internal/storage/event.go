package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventTitle         = errors.New("empty event title")
	ErrEventDate          = errors.New("empty event date")
	ErrEventEndDate       = errors.New("empty event end date")
	ErrEventUserID        = errors.New("bad event user id")
	ErrEventID            = errors.New("bad event id")
	ErrNotExist           = errors.New("event not exist")
	ErrNotValidRemindDate = errors.New("remind date invalid")
	ErrNotValidEventDate  = errors.New("date is more then end date")
	ErrDuplicateEvent     = errors.New("duplicate event error")
)

type Event struct {
	ID         string    `db:"id" json:"id"`
	Title      string    `db:"title" json:"title"`
	Date       time.Time `db:"date" json:"date"`
	EndDate    time.Time `db:"end_date" json:"end_date"` //nolint:tagliatelle
	Desc       string    `db:"description" json:"description"`
	UserID     string    `db:"user_id" json:"user_id"`         //nolint:tagliatelle
	RemindDate time.Time `db:"remind_date" json:"remind_date"` //nolint:tagliatelle
}

type EventDTO struct {
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	EndDate    time.Time `json:"end_date"` //nolint:tagliatelle
	Desc       string    `json:"description"`
	UserID     string    `json:"user_id"`     //nolint:tagliatelle
	RemindDate time.Time `json:"remind_date"` //nolint:tagliatelle
}

func (e *EventDTO) Transfer() Event {
	return Event{
		Title:      e.Title,
		Date:       e.Date,
		EndDate:    e.EndDate,
		Desc:       e.Desc,
		UserID:     e.UserID,
		RemindDate: e.RemindDate,
	}
}

type EventUseCase struct {
	storage Storage
}

func (u *EventUseCase) validateEvent(e Event) error {
	switch {
	case e.Title == "":
		return ErrEventTitle
	case e.EndDate.IsZero():
		return ErrEventEndDate
	case e.Date.IsZero():
		return ErrEventDate
	case !e.RemindDate.IsZero() && (!e.RemindDate.After(e.Date) || !e.RemindDate.Before(e.EndDate)):
		return ErrNotValidRemindDate
	}

	if e.Date.After(e.EndDate) {
		return ErrNotValidEventDate
	}

	if _, err := uuid.Parse(e.UserID); err != nil {
		return ErrEventUserID
	}

	return nil
}

func (u *EventUseCase) getAndUpdateEventValue(ctx context.Context, eventID string, modifyEvent Event) (Event, error) {
	event, err := u.storage.GetEvent(ctx, eventID)
	if err != nil {
		return Event{}, err
	}
	if modifyEvent.Title != "" {
		event.Title = modifyEvent.Title
	}
	if !modifyEvent.Date.IsZero() {
		event.Date = modifyEvent.Date
	}
	if !modifyEvent.EndDate.IsZero() {
		event.EndDate = modifyEvent.EndDate
	}
	if !modifyEvent.RemindDate.IsZero() {
		event.RemindDate = modifyEvent.RemindDate
	}
	if modifyEvent.Desc != "" {
		event.Desc = modifyEvent.Desc
	}
	if modifyEvent.UserID != "" {
		event.UserID = modifyEvent.UserID
	}
	return event, nil
}

func (u *EventUseCase) Create(ctx context.Context, event Event) (string, error) {
	err := u.validateEvent(event)
	if err != nil {
		return "", fmt.Errorf("EventUseCase - CreateEvent - u.storage.CreateEvent: %w", err)
	}
	res, err := u.storage.CreateEvent(ctx, event)
	if err != nil {
		return "", fmt.Errorf("EventUseCase - CreateEvent - u.storage.CreateEvent: %w", err)
	}
	return res, nil
}

func (u *EventUseCase) Update(ctx context.Context, eventID string, event Event) error {
	if _, err := uuid.Parse(eventID); err != nil {
		return ErrEventID
	}
	mEvent, err := u.getAndUpdateEventValue(ctx, eventID, event)
	if err != nil {
		return fmt.Errorf("EventUseCase - UpdateEvent - u.storage.UpdateEvent: %w", err)
	}
	err = u.validateEvent(mEvent)
	if err != nil {
		return fmt.Errorf("EventUseCase - UpdateEvent - u.storage.UpdateEvent: %w", err)
	}
	err = u.storage.UpdateEvent(ctx, eventID, event)
	if err != nil {
		return fmt.Errorf("EventUseCase - UpdateEvent - u.storage.UpdateEvent: %w", err)
	}
	return nil
}

func (u *EventUseCase) Delete(ctx context.Context, eventID string) error {
	if _, err := uuid.Parse(eventID); err != nil {
		return ErrEventID
	}
	err := u.storage.DeleteEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("EventUseCase - DeleteEvent - u.storage.DeleteEvent: %w", err)
	}
	return nil
}

func (u *EventUseCase) DeleteBeforeDate(ctx context.Context, date time.Time) error {
	err := u.storage.DeleteEventsBeforeDate(ctx, date)
	if err != nil {
		return fmt.Errorf("EventUseCase - DeleteBeforeDate - u.storage.DeleteEventsBeforeDate: %w", err)
	}
	return nil
}

func (u *EventUseCase) GetDailyEvents(ctx context.Context, date time.Time) ([]Event, error) {
	events, err := u.storage.GetDailyEvents(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - GetDailyEvents - u.storage.GetDailyEvents: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) GetWeeklyEvents(ctx context.Context, date time.Time) ([]Event, error) {
	events, err := u.storage.GetWeeklyEvents(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - GetWeeklyEvents - u.storage.GetWeeklyEvents: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) GetMonthlyEvents(ctx context.Context, date time.Time) ([]Event, error) {
	events, err := u.storage.GetMonthlyEvents(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - MonthlyEvents - u.storage.MonthlyEvents: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) GetEventsByPeriod(ctx context.Context, start, end time.Time) ([]Event, error) {
	events, err := u.storage.GetEventsByPeriod(ctx, start, end)
	if err != nil {
		return nil, fmt.Errorf("EventUseCase - GetEventsByPeriod - u.storage.GetEventsByPeriod: %w", err)
	}
	return events, nil
}

func (u *EventUseCase) Connect(ctx context.Context) error {
	err := u.storage.Connect(ctx)
	if err != nil {
		return fmt.Errorf("EventUseCase - Connect - u.storage.Connect: %w", err)
	}
	return nil
}

func (u *EventUseCase) Close(ctx context.Context) error {
	err := u.storage.Close(ctx)
	if err != nil {
		return fmt.Errorf("EventUseCase - Connect - u.storage.Close: %w", err)
	}
	return nil
}

func New(st Storage) *EventUseCase {
	return &EventUseCase{
		storage: st,
	}
}

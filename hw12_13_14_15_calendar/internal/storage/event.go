package storage

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrEventTitle    = errors.New("empty event title")
	ErrEventDate     = errors.New("empty event date")
	ErrEventDuration = errors.New("empty event duration")
	ErrEventUserID   = errors.New("empty event user id")
	ErrNotExist      = errors.New("event not exist")
)

type Event struct {
	ID         int32         `db:"id"`
	Title      string        `db:"title"`
	Date       time.Time     `db:"date"`
	Duration   time.Duration `db:"duration"`
	Desc       string        `db:"desc"`
	UserID     int32         `db:"user_id"`
	RemindDate time.Duration `db:"remind_date"`
}

type Storage interface {
	CreateEvent(context.Context, Event) error
	UpdateEvent(context.Context, Event) error
	DeleteEvent(context.Context, int32) error
	GetDailyEvents(context.Context, time.Time) ([]Event, error)
	GetWeeklyEvents(context.Context, time.Time) ([]Event, error)
	GetMonthlyEvents(context.Context, time.Time) ([]Event, error)
	validateEvent(Event) error
}

type EventUseCase struct {
	storage Storage
}

func (u *EventUseCase) Create(ctx context.Context, event Event) error {
	err := u.storage.CreateEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("EventUseCase - CreateEvent - u.storage.CreateEvent: %w", err)
	}
	return nil
}

func (u *EventUseCase) Update(ctx context.Context, event Event) error {
	err := u.storage.UpdateEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("EventUseCase - UpdateEvent - u.storage.UpdateEvent: %w", err)
	}
	return nil
}

func (u *EventUseCase) Delete(ctx context.Context, eventID int32) error {
	err := u.storage.DeleteEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("EventUseCase - DeleteEvent - u.storage.DeleteEvent: %w", err)
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

func New(st Storage) *EventUseCase {
	return &EventUseCase{
		storage: st,
	}
}

func ValidateEvent(event Event) error {
	switch {
	case event.Title == "":
		return ErrEventTitle
	case event.UserID == 0:
		return ErrEventUserID
	case event.Date.IsZero():
		return ErrEventDate
	case event.Duration == 0:
		return ErrEventDuration
	}
	return nil
}

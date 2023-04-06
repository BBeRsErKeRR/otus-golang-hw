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
	ErrEventUserID        = errors.New("empty event user id")
	ErrNotExist           = errors.New("event not exist")
	ErrNotValidRemindDate = errors.New("remind date invalid")
	ErrNotValidEventDate  = errors.New("date is more then end date")
)

type Event struct {
	ID         string    `db:"id"`
	Title      string    `db:"title"`
	Date       time.Time `db:"date"`
	EndDate    time.Time `db:"end_date"`
	Desc       string    `db:"description"`
	UserID     string    `db:"user_id"`
	RemindDate time.Time `db:"remind_date"`
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

func (u *EventUseCase) Update(ctx context.Context, eventID interface{}, event Event) error {
	err := u.storage.UpdateEvent(ctx, eventID, event)
	if err != nil {
		return fmt.Errorf("EventUseCase - UpdateEvent - u.storage.UpdateEvent: %w", err)
	}
	return nil
}

func (u *EventUseCase) Delete(ctx context.Context, eventID interface{}) error {
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

func ValidateEvent(e Event) error {
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

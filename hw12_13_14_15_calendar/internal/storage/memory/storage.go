package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func (st *Storage) GetEvent(ctx context.Context, eventID string) (storage.Event, error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	event, ok := st.events[eventID]
	if !ok {
		return storage.Event{}, storage.ErrNotExist
	}
	return event, nil
}

func (st *Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	for _, e := range st.events {
		if e.Date.Equal(event.Date) && e.EndDate.Equal(event.EndDate) && e.Title == event.Title {
			return "", storage.ErrDuplicateEvent
		}
	}
	event.ID = uuid.New().String()
	st.events[event.ID] = event
	return event.ID, nil
}

func (st *Storage) UpdateEvent(ctx context.Context, eventID string, modifyEvent storage.Event) error {
	sEvent, err := st.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	if modifyEvent.Title == "" {
		modifyEvent.Title = sEvent.Title
	}
	if modifyEvent.Date.IsZero() {
		modifyEvent.Date = sEvent.Date
	}
	if modifyEvent.EndDate.IsZero() {
		modifyEvent.EndDate = sEvent.EndDate
	}
	if modifyEvent.RemindDate.IsZero() {
		modifyEvent.RemindDate = sEvent.RemindDate
	}
	if modifyEvent.Desc == "" {
		modifyEvent.Desc = sEvent.Desc
	}
	if modifyEvent.UserID == "" {
		modifyEvent.UserID = sEvent.UserID
	}
	st.events[eventID] = modifyEvent
	return nil
}

func (st *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.events[eventID]; !ok {
		return storage.ErrNotExist
	}
	delete(st.events, eventID)
	return nil
}

func (st *Storage) DeleteEventsBeforeDate(ctx context.Context, date time.Time) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	for id, e := range st.events {
		if e.EndDate.After(date) {
			delete(st.events, id)
		}
	}
	return nil
}

func (st *Storage) GetEventsByPeriod(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	res := make([]storage.Event, 0, len(st.events))
	for _, e := range st.events {
		if (e.Date.After(start) && e.Date.Before(end)) || (e.Date.Before(start) && e.EndDate.After(end)) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (st *Storage) GetDailyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.GetEventsByPeriod(ctx, date, date.AddDate(0, 0, 1))
}

func (st *Storage) GetWeeklyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.GetEventsByPeriod(ctx, date, date.AddDate(0, 0, 7))
}

func (st *Storage) GetMonthlyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.GetEventsByPeriod(ctx, date, date.AddDate(0, 1, 0))
}

func (st *Storage) Connect(ctx context.Context) error {
	return nil
}

func (st *Storage) Close(ctx context.Context) error {
	return nil
}

func New() *Storage {
	return &Storage{
		events: map[string]storage.Event{},
	}
}

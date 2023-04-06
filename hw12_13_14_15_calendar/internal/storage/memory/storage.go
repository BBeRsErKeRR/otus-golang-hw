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

func (st *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	err := st.ValidateEvent(event)
	if err != nil {
		return err
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	event.ID = uuid.New().String()
	st.events[event.ID] = event
	return nil
}

func (st *Storage) UpdateEvent(ctx context.Context, eventID string, event storage.Event) error {
	err := st.ValidateEvent(event)
	if err != nil {
		return err
	}

	st.mu.Lock()
	defer st.mu.Unlock()
	_, ok := st.events[eventID]
	if !ok {
		return storage.ErrNotExist
	}

	st.events[event.ID] = event
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

func (st *Storage) getEventsByPeriod(start, end time.Time) ([]storage.Event, error) {
	res := make([]storage.Event, 0, len(st.events))
	for _, e := range st.events {
		if e.Date.After(start) && e.EndDate.Before(end) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (st *Storage) GetDailyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.getEventsByPeriod(date, date.AddDate(0, 0, 1))
}

func (st *Storage) GetWeeklyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.getEventsByPeriod(date, date.AddDate(0, 0, 7))
}

func (st *Storage) GetMonthlyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.getEventsByPeriod(date, date.AddDate(0, 1, 0))
}

func (st *Storage) ValidateEvent(event storage.Event) error {
	return storage.ValidateEvent(event)
}

func New() *Storage {
	return &Storage{
		events: map[string]storage.Event{},
	}
}

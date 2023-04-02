package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[int32]storage.Event
	mu     sync.RWMutex
}

func (st *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	err := st.eventValidate(event)
	if err != nil {
		return err
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	event.ID = int32(uuid.New().ID())
	st.events[event.ID] = event
	return nil
}

func (st *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	err := st.eventValidate(event)
	if err != nil {
		return err
	}

	st.mu.Lock()
	defer st.mu.Unlock()
	_, ok := st.events[event.ID]
	if !ok {
		return storage.ErrNotExist
	}

	st.events[event.ID] = event
	return nil
}

func (st *Storage) DeleteEvent(ctx context.Context, eventId int32) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.events[eventId]; !ok {
		return storage.ErrNotExist
	}
	delete(st.events, eventId)
	return nil
}

func (st *Storage) getEventsByPeriod(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	res := make([]storage.Event, 0, len(st.events))
	for _, e := range st.events {
		if e.Date.After(start) && e.Date.Before(end) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (st *Storage) GetDailyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.getEventsByPeriod(ctx, date, date.AddDate(0, 0, 1))
}

func (st *Storage) GetWeeklyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.getEventsByPeriod(ctx, date, date.AddDate(0, 0, 7))
}

func (st *Storage) GetMonthlyEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return st.getEventsByPeriod(ctx, date, date.AddDate(0, 1, 0))
}

func (st *Storage) eventValidate(event storage.Event) error {
	return storage.ValidateEvent(event)
}

func New() *Storage {
	return &Storage{
		events: map[int32]storage.Event{},
	}
}

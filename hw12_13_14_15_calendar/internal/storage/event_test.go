package storage

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type tStorage struct{}

func (ts *tStorage) CreateEvent(ctx context.Context, event Event) error {
	return nil
}

func (ts *tStorage) UpdateEvent(ctx context.Context, id string, event Event) error {
	return nil
}

func (ts *tStorage) DeleteEvent(ctx context.Context, id string) error {
	return nil
}

func (ts *tStorage) GetDailyEvents(ctx context.Context, date time.Time) ([]Event, error) {
	return []Event{}, nil
}

func (ts *tStorage) GetWeeklyEvents(ctx context.Context, date time.Time) ([]Event, error) {
	return []Event{}, nil
}

func (ts *tStorage) GetMonthlyEvents(ctx context.Context, date time.Time) ([]Event, error) {
	return []Event{}, nil
}

func (ts *tStorage) Connect(ctx context.Context) error {
	return nil
}

func (ts *tStorage) Close(ctx context.Context) error {
	return nil
}

func TestStorage(t *testing.T) {
	ctx := context.Background()
	testcases := []struct {
		Name   string
		Event  Event
		Action func(ctx context.Context, u *EventUseCase, event Event) error
		Error  error
	}{
		{
			Name:  "invalid title",
			Event: Event{UserID: uuid.New().String(), Date: time.Now(), EndDate: time.Now().Add(4 * time.Hour)},
			Action: func(ctx context.Context, u *EventUseCase, event Event) error {
				return u.Create(context.Background(), event)
			},
			Error: ErrEventTitle,
		},
		{
			Name:  "invalid end date",
			Event: Event{Title: "some event 1", UserID: uuid.New().String(), Date: time.Now()},
			Action: func(ctx context.Context, u *EventUseCase, event Event) error {
				return u.Create(context.Background(), event)
			},
			Error: ErrEventEndDate,
		},
		{
			Name:  "invalid date",
			Event: Event{Title: "some event 1", UserID: uuid.New().String(), EndDate: time.Now().Add(22 * time.Hour)},
			Action: func(ctx context.Context, u *EventUseCase, event Event) error {
				return u.Update(context.Background(), "", event)
			},
			Error: ErrEventDate,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			u := New(&tStorage{})
			err := testcase.Action(ctx, u, testcase.Event)
			require.ErrorIs(t, err, testcase.Error)
		})
	}
}

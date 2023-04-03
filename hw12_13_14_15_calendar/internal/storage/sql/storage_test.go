package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var errEmpty = errors.New("empty result")

func TestStorage(t *testing.T) {
	ctx := context.Background()
	config := storage.Config{
		Host:     "localhost",
		Port:     "5532",
		Storage:  "sql",
		Driver:   "postgres",
		Ssl:      "disable",
		Database: "calendar",
		User:     "calendar",
		Password: "passwd",
	}

	testcases := []struct {
		Name   string
		Event  storage.Event
		Action func(event storage.Event, st *Storage, ctx context.Context) error
		Err    error
	}{
		{
			Name: "check crud",
			Event: storage.Event{
				Title:      "some event 1",
				Desc:       "this is the test event 1",
				UserID:     uuid.New().String(),
				Date:       time.Now(),
				Duration:   2 * time.Hour,
				RemindDate: time.Hour,
			},
			Action: func(event storage.Event, st *Storage, ctx context.Context) error {
				err := st.CreateEvent(ctx, event)
				if err != nil {
					return err
				}
				events, err := st.GetDailyEvents(ctx, time.Now().Add(-time.Hour))
				if err != nil {
					return err
				}
				if len(events) == 0 {
					return errEmpty
				}
				newTitle := "modified event 1"
				newDate := time.Now().AddDate(0, 0, 1)

				mEvent := events[0]
				mEvent.Title = newTitle
				mEvent.Date = newDate
				err = st.UpdateEvent(ctx, mEvent.ID, mEvent)
				if err != nil {
					return err
				}
				mEvents, err := st.GetDailyEvents(ctx, newDate.Add(-time.Hour))
				if err != nil {
					return err
				}
				if len(mEvents) == 0 {
					return errEmpty
				}
				assertEvent := mEvents[0]
				if assertEvent.ID != mEvent.ID ||
					assertEvent.Title != newTitle ||
					assertEvent.Date != newDate {
					return fmt.Errorf("failed assertion got %v found %v", event, assertEvent)
				}
				return nil
			},
		},
		{
			Name: "invalid title",
			Event: storage.Event{
				Desc:       "this is the test event 1",
				UserID:     uuid.New().String(),
				Date:       time.Now(),
				Duration:   2 * time.Hour,
				RemindDate: time.Hour,
			},
			Action: func(event storage.Event, st *Storage, ctx context.Context) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventTitle,
		},
		{
			Name: "invalid duration",
			Event: storage.Event{
				Title:      "some event 1",
				Desc:       "this is the test event 1",
				UserID:     uuid.New().String(),
				Date:       time.Now(),
				RemindDate: time.Hour,
			},
			Action: func(event storage.Event, st *Storage, ctx context.Context) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventDuration,
		},
		{
			Name: "invalid date",
			Event: storage.Event{
				Title:      "some event 1",
				Desc:       "this is the test event 1",
				UserID:     uuid.New().String(),
				Duration:   2 * time.Hour,
				RemindDate: time.Hour,
			},
			Action: func(event storage.Event, st *Storage, ctx context.Context) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventDate,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			st := New(&config)
			require.NoError(t, st.Connect(ctx))
			err := testcase.Action(testcase.Event, st, ctx)
			require.NoError(t, st.Close(ctx))
			require.ErrorIs(t, err, testcase.Err)
		})
	}
}

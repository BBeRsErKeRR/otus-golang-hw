package memorystorage

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
	testcases := []struct {
		Name   string
		Event  storage.Event
		Action func(ctx context.Context, st *Storage, event storage.Event) error
		Err    error
	}{
		{
			Name: "check crud",
			Event: storage.Event{
				Title:      "some event 1",
				Desc:       "this is the test event 1",
				UserID:     uuid.New().String(),
				Date:       time.Now(),
				EndDate:    time.Now().Add(4 * time.Hour),
				RemindDate: time.Now().Add(3 * time.Hour),
			},
			Action: func(ctx context.Context, st *Storage, event storage.Event) error {
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
				mEvent.EndDate = newDate.Add(4 * time.Hour)
				mEvent.RemindDate = newDate.Add(3 * time.Hour)
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
				UserID:  uuid.New().String(),
				Date:    time.Now(),
				EndDate: time.Now().Add(4 * time.Hour),
			},
			Action: func(ctx context.Context, st *Storage, event storage.Event) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventTitle,
		},
		{
			Name: "invalid duration",
			Event: storage.Event{
				Title:  "some event 1",
				UserID: uuid.New().String(),
				Date:   time.Now(),
			},
			Action: func(ctx context.Context, st *Storage, event storage.Event) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventEndDate,
		},
		{
			Name: "invalid date",
			Event: storage.Event{
				Title:   "some event 1",
				UserID:  uuid.New().String(),
				EndDate: time.Now().Add(4 * time.Hour),
			},
			Action: func(ctx context.Context, st *Storage, event storage.Event) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventDate,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			memStorage := New()
			err := testcase.Action(ctx, memStorage, testcase.Event)
			require.ErrorIs(t, err, testcase.Err)
		})
	}

	t.Run("some get check", func(t *testing.T) {
		memStorage := Storage{
			events: map[string]storage.Event{
				uuid.New().String(): {
					Title:      "not exist",
					UserID:     uuid.New().String(),
					Date:       time.Now().Add(-48 * time.Hour),
					EndDate:    time.Now().Add(-24 * time.Hour),
					RemindDate: time.Now().Add(-26 * time.Hour),
				},
				uuid.New().String(): {
					Title:      "some event 1",
					Desc:       "this is the test event 1",
					UserID:     uuid.New().String(),
					Date:       time.Now(),
					EndDate:    time.Now().Add(4 * time.Hour),
					RemindDate: time.Now().Add(3 * time.Hour),
				},
				uuid.New().String(): {
					Title:      "some event 2",
					Desc:       "this is the test event 2",
					UserID:     uuid.New().String(),
					Date:       time.Now().Add(24 * time.Hour),
					EndDate:    time.Now().Add(26 * time.Hour),
					RemindDate: time.Now().Add(25 * time.Hour),
				},
				uuid.New().String(): {
					Title:      "some event 3",
					Desc:       "this is the test event 3",
					UserID:     uuid.New().String(),
					Date:       time.Now().AddDate(0, 0, 14),
					EndDate:    time.Now().AddDate(0, 0, 15),
					RemindDate: time.Now().AddDate(0, 0, 14).Add(3 * time.Hour),
				},
			},
		}
		testDate := time.Now().Add(-4 * time.Minute)
		dayEvents, err := memStorage.GetDailyEvents(ctx, testDate)
		require.NoError(t, err)
		require.Equal(t, 1, len(dayEvents))
		weekEvents, err := memStorage.GetWeeklyEvents(ctx, testDate)
		require.NoError(t, err)
		require.Equal(t, 2, len(weekEvents))
		monthEvents, err := memStorage.GetMonthlyEvents(ctx, testDate)
		require.NoError(t, err)
		require.Equal(t, 3, len(monthEvents))
	})
}

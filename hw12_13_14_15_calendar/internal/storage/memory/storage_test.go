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

func TestStorage(t *testing.T) { //nolint:gocognit
	now := time.Now()
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
				Title: "some event 1", Desc: "this is the test event 1", UserID: uuid.New().String(),
				Date: now, EndDate: now.Add(4 * time.Hour), RemindDate: now.Add(-3 * time.Hour),
			},
			Action: func(ctx context.Context, st *Storage, event storage.Event) error {
				_, err := st.CreateEvent(ctx, event)
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
			Name: "check duplicate error",
			Event: storage.Event{
				Title: "some event 2", Desc: "this is the test event 2", UserID: uuid.New().String(),
				Date: now, EndDate: now.Add(4 * time.Hour), RemindDate: now.Add(-3 * time.Hour),
			},
			Err: storage.ErrDuplicateEvent,
			Action: func(ctx context.Context, st *Storage, event storage.Event) error {
				id, err := st.CreateEvent(ctx, event)
				if err != nil {
					return err
				}
				_, err = st.GetEvent(ctx, id)
				if err != nil {
					return err
				}
				_, err = st.CreateEvent(ctx, event)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "check kindReminder",
			Event: storage.Event{
				Title: "some event 3", Desc: "this is the test event 3", UserID: uuid.New().String(),
				Date: now, EndDate: now.Add(4 * time.Hour), RemindDate: now.Add(-3 * time.Minute),
			},
			Action: func(ctx context.Context, st *Storage, event storage.Event) error {
				_, err := st.CreateEvent(ctx, event)
				if err != nil {
					return err
				}
				events, err := st.GetKindReminder(ctx, now.Add(-5*time.Second))
				if err != nil {
					return err
				}
				if len(events) == 0 {
					return errEmpty
				}
				return nil
			},
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
		title := "not exist"
		date := now.Add(-48 * time.Hour)
		endDate := now.Add(-24 * time.Hour)
		testDate := now.Add(-4 * time.Minute)
		memStorage := Storage{
			events: map[string]storage.Event{
				uuid.New().String(): {
					Title: title, UserID: uuid.New().String(),
					Date: date, EndDate: endDate, RemindDate: now.Add(-26 * time.Hour),
				},
				uuid.New().String(): {
					Title: "te2", Desc: "te2", UserID: uuid.New().String(),
					Date: now, EndDate: now.Add(4 * time.Hour), RemindDate: now.Add(3 * time.Hour),
				},
				uuid.New().String(): {
					Title: "some event 2", Desc: "this is the test event 2", UserID: uuid.New().String(),
					Date: now.Add(24 * time.Hour), EndDate: now.Add(26 * time.Hour), RemindDate: now.Add(25 * time.Hour),
				},
				uuid.New().String(): {
					Title: "some event 3", Desc: "this is the test event 3", UserID: uuid.New().String(),
					Date: now.AddDate(0, 0, 14), EndDate: now.AddDate(0, 0, 15), RemindDate: now.AddDate(0, 0, 14).Add(3 * time.Hour),
				},
			},
		}
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

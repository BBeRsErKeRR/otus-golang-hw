//go:build integration

package sqlstorage

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var errEmpty = errors.New("empty result")

//go:embed fixture/events.sql
var query string

func getUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

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
		Action func(ctx context.Context, st *Storage, event storage.Event, db *sqlx.DB) error
		Err    error
	}{
		{
			Name: "check crud",
			Event: storage.Event{
				Title:      "some event 1",
				Desc:       "this is the test event 1",
				UserID:     uuid.New().String(),
				Date:       time.Now(),
				EndDate:    time.Now().Add(3 * time.Hour),
				RemindDate: time.Now().Add(2 * time.Hour),
			},
			Action: func(ctx context.Context, st *Storage, event storage.Event, db *sqlx.DB) error {
				if _, err := db.Exec(`TRUNCATE TABLE events`); err != nil {
					return err
				}
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
				if assertEvent.ID != mEvent.ID {
					return fmt.Errorf("failed assertion expect '%v' found '%v'", mEvent.ID, assertEvent.ID)
				}
				if assertEvent.Title != newTitle {
					return fmt.Errorf("failed title assertion expect '%v' found '%v", newTitle, assertEvent.Title)
				}
				if newDate.Equal(assertEvent.Date) {
					return fmt.Errorf("failed date assertion expect '%v' found '%v", newDate, assertEvent.Date)
				}
				return nil
			},
		},
		{
			Name:  "invalid title",
			Event: storage.Event{UserID: uuid.New().String(), Date: time.Now(), EndDate: time.Now().Add(4 * time.Hour)},
			Action: func(ctx context.Context, st *Storage, event storage.Event, db *sqlx.DB) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventTitle,
		},
		{
			Name:  "invalid end date",
			Event: storage.Event{Title: "some event 1", UserID: uuid.New().String(), Date: time.Now()},
			Action: func(ctx context.Context, st *Storage, event storage.Event, db *sqlx.DB) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventEndDate,
		},
		{
			Name:  "invalid date",
			Event: storage.Event{Title: "some event 1", UserID: uuid.New().String(), EndDate: time.Now().Add(22 * time.Hour)},
			Action: func(ctx context.Context, st *Storage, event storage.Event, db *sqlx.DB) error {
				return st.CreateEvent(context.Background(), event)
			},
			Err: storage.ErrEventDate,
		},
		{
			Name: "some get check",
			Action: func(ctx context.Context, st *Storage, event storage.Event, db *sqlx.DB) error {
				if _, err := db.Exec(query); err != nil {
					return err
				}
				testDate := time.Now().Add(-4 * time.Minute)
				dayEvents, err := st.GetDailyEvents(ctx, testDate)
				if err != nil {
					return err
				}
				if len(dayEvents) != 1 {
					return fmt.Errorf("failed check date '%v' dayEvents '%v' found '%v'", testDate, 1, dayEvents)
				}
				weekEvents, err := st.GetWeeklyEvents(ctx, testDate)
				if err != nil {
					return err
				}
				if len(weekEvents) != 2 {
					return fmt.Errorf("failed check date '%v' weekEvents '%v' found '%v'", testDate, 2, weekEvents)
				}
				monthEvents, err := st.GetMonthlyEvents(ctx, testDate)
				if err != nil {
					return err
				}
				if len(monthEvents) != 3 {
					return fmt.Errorf("failed check date '%v' monthEvents '%v' found '%v'", testDate, 3, monthEvents)
				}
				return nil
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			st := New(&config)
			require.NoError(t, st.Connect(ctx))
			field := reflect.Indirect(reflect.ValueOf(st)).FieldByName("db")
			db := getUnexportedField(field).(*sqlx.DB)
			err := testcase.Action(ctx, st, testcase.Event, db)
			require.NoError(t, st.Close(ctx))
			require.ErrorIs(t, err, testcase.Err)
		})
	}
}

package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib" //nolint:blank-imports
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotInitDB  = errors.New("not init db")
	ErrFailUpdate = errors.New("update failed")
)

type Storage struct {
	dataSourceName string
	db             *sqlx.DB
}

func (st *Storage) Connect(ctx context.Context) error {
	var err error
	st.db, err = sqlx.ConnectContext(ctx, "pgx", st.dataSourceName)
	if err != nil {
		return err
	}
	err = st.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (st *Storage) Close(ctx context.Context) error {
	if st.db == nil {
		return ErrNotInitDB
	}
	return st.db.Close()
}

func New(conf *storage.Config) *Storage {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Ssl,
	)

	return &Storage{
		dataSourceName: dsn,
	}
}

func (st *Storage) execNamedQuery(ctx context.Context, query string, event storage.Event) (sql.Result, error) {
	if st.db == nil {
		return nil, ErrNotInitDB
	}
	stmt, err := st.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("can't prepare query:\n%v\n%w", query, err)
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, event)
	if err != nil {
		return nil, err
	}

	return res, err
}

const createEventQ = `
INSERT INTO events(id, title, date, end_date, description, user_id, remind_date) 
	VALUES (:id, :title, :date, :end_date, :description, :user_id, :remind_date)
`

func (st *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	err := st.ValidateEvent(event)
	if err != nil {
		return err
	}
	event.ID = uuid.New().String()
	_, err = st.execNamedQuery(ctx, createEventQ, event)
	return err
}

const updateEventQ = `
UPDATE events 
	SET title = :title, 
		date = :date, 
		end_date = :end_date, 
		description = :description, 
		user_id = :user_id,
		remind_date = :remind_date
	WHERE id = :id
`

func (st *Storage) UpdateEvent(ctx context.Context, eventID string, event storage.Event) error {
	err := st.ValidateEvent(event)
	if err != nil {
		return err
	}
	event.ID = eventID

	res, err := st.execNamedQuery(ctx, updateEventQ, event)
	if err != nil {
		return err
	}
	rowA, err := res.RowsAffected()
	if rowA == 0 {
		return ErrFailUpdate
	}
	return err
}

const deleteEventQ = `DELETE FROM events WHERE id=$1`

func (st *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	if st.db == nil {
		return ErrNotInitDB
	}
	_, err := st.db.ExecContext(ctx, deleteEventQ, eventID)
	return err
}

const getEventsByPeriodQ = `
SELECT * FROM events 
WHERE date>=$1 
	AND end_date<=$2
`

func (st *Storage) getEventsByPeriod(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	var res []storage.Event
	if st.db == nil {
		return res, ErrNotInitDB
	}
	err := st.db.SelectContext(ctx, &res, getEventsByPeriodQ, start.UTC(), end.UTC())
	if err != nil {
		return res, err
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

func (st *Storage) ValidateEvent(event storage.Event) error {
	return storage.ValidateEvent(event)
}

package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" //nolint:blank-imports
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
	fmt.Println(st.dataSourceName)
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

const getEventQ = `
SELECT * FROM events 
WHERE id=$1
`

func (st *Storage) GetEvent(ctx context.Context, eventID string) (storage.Event, error) {
	var res storage.Event
	if st.db == nil {
		return res, ErrNotInitDB
	}
	row := st.db.QueryRowxContext(ctx, getEventQ, eventID)
	err := row.StructScan(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

const createEventQ = `
INSERT INTO events(id, title, date, end_date, description, user_id, remind_date) 
	VALUES (:id, :title, :date, :end_date, :description, :user_id, :remind_date)
`

func (st *Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	event.ID = uuid.New().String()
	_, err := st.execNamedQuery(ctx, createEventQ, event)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return "", storage.ErrDuplicateEvent
		}
		return "", err
	}
	return event.ID, nil
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

func (st *Storage) UpdateEvent(ctx context.Context, eventID string, modifyEvent storage.Event) error {
	sEvent, err := st.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}
	modifyEvent.ID = eventID
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
	res, err := st.execNamedQuery(ctx, updateEventQ, modifyEvent)
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

const deleteBeforeEventQ = `DELETE FROM events WHERE end_date <= $1`

func (st *Storage) DeleteEventsBeforeDate(ctx context.Context, date time.Time) error {
	if st.db == nil {
		return ErrNotInitDB
	}
	_, err := st.db.ExecContext(ctx, deleteBeforeEventQ, date)
	return err
}

const getEventsByPeriodQ = `
SELECT * FROM events 
WHERE (date>=$1 AND date<=$2)
	OR (date<$1 AND end_date>$2)
`

func (st *Storage) GetEventsByPeriod(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	var res []storage.Event
	if st.db == nil {
		return res, ErrNotInitDB
	}
	err := st.db.SelectContext(ctx, &res, getEventsByPeriodQ, start, end)
	if err != nil {
		return res, err
	}
	return res, nil
}

const getReminderEvents = `
SELECT * FROM events
WHERE date>=$1 AND remind_date<=$1
`

func (st *Storage) GetKindReminder(ctx context.Context, date time.Time) ([]storage.Event, error) {
	var res []storage.Event
	if st.db == nil {
		return res, ErrNotInitDB
	}
	err := st.db.SelectContext(ctx, &res, getReminderEvents, date)
	if err != nil {
		return res, err
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

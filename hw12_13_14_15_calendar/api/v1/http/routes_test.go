package v1routes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestV1HTTPHandlers(t *testing.T) {
	loggerConfig := logger.Config{
		Level:    "debug",
		OutPaths: []string{},
		ErrPaths: []string{},
	}
	logger, err := logger.New(&loggerConfig)
	require.NoError(t, err)
	st := memorystorage.New()
	application := app.New(logger, *storage.New(st))
	router := mux.NewRouter()
	handler := NewHandler(application, logger)
	handler.AddV1Routes(router)

	testcases := map[string]struct {
		URL         func(t *testing.T, app *app.App) string
		Method      string
		Body        string
		Action      http.HandlerFunc
		CheckResult func(t *testing.T, res *httptest.ResponseRecorder, req *http.Request, st storage.Storage)
	}{
		"check createEvent": {
			URL:    func(t *testing.T, app *app.App) string { return "/v1/event" }, //nolint:thelper
			Method: http.MethodPost,
			Body: `{
					"title": "test1",
					"user_id": "2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
					"date": "2023-01-01T16:00:00Z",
					"end_date": "2023-01-02T16:00:00Z"
				  }`,
			Action: handler.CreateEvent,
			CheckResult: func(t *testing.T, res *httptest.ResponseRecorder, req *http.Request, st storage.Storage) { //nolint:thelper,lll
				var data map[string]string
				err := json.Unmarshal(res.Body.Bytes(), &data)
				require.NoError(t, err)
				require.Equal(t, "Created", data["msg"], fmt.Sprintf("expected 'msg' with 'Created', actual: %v", data))
			},
		},
		"check deleteEvent": {
			URL: func(t *testing.T, app *app.App) string { //nolint:thelper
				event := storage.Event{
					Title:   "Deleted",
					Date:    time.Now(),
					EndDate: time.Now().AddDate(0, 0, 1),
					UserID:  uuid.New().String(),
				}
				id, err := app.CreateEvent(context.Background(), event)
				require.NoError(t, err)
				return fmt.Sprintf("/v1/event/%v", id)
			},
			Method: http.MethodDelete,
			Action: handler.DeleteEvent,
			CheckResult: func(t *testing.T, res *httptest.ResponseRecorder, req *http.Request, st storage.Storage) { //nolint:thelper,lll
				var data map[string]string
				err := json.Unmarshal(res.Body.Bytes(), &data)
				require.NoError(t, err)
				require.Equal(t, "Deleted", data["msg"])
			},
		},
		"check updateEvent": {
			URL: func(t *testing.T, app *app.App) string { //nolint:thelper
				event := storage.Event{
					Title:   "Updated",
					Date:    time.Now(),
					EndDate: time.Now().AddDate(0, 0, 1),
					UserID:  uuid.New().String(),
				}
				id, err := app.CreateEvent(context.Background(), event)
				require.NoError(t, err)
				return fmt.Sprintf("/v1/event/%v", id)
			},
			Body: `{
				"title": "test2",
				"user_id": "2bb0d64e-8f6e-4863-b1d8-8b20018c743d",
				"date": "2023-01-01T16:00:00Z",
				"end_date": "2023-01-02T16:00:00Z"
			  }`,
			Method: http.MethodPut,
			Action: handler.UpdateEvent,
			CheckResult: func(t *testing.T, res *httptest.ResponseRecorder, req *http.Request, st storage.Storage) { //nolint:thelper,lll
				var data map[string]string
				err := json.Unmarshal(res.Body.Bytes(), &data)
				require.NoError(t, err)
				require.Equal(t, "Updated", data["msg"], fmt.Sprintf("expected Updated, found %v", data))
				id := path.Base(req.URL.Path)
				event, err := st.GetEvent(context.Background(), id)
				require.NoError(t, err)
				require.Equal(t, "test2", event.Title, fmt.Sprintf("expected test2, found %v", event))
			},
		},
		"check GetDaily": {
			URL: func(t *testing.T, app *app.App) string { //nolint:thelper
				event := storage.Event{
					Title:   "GetDaily",
					Date:    time.Now(),
					EndDate: time.Now().AddDate(0, 0, 1),
					UserID:  uuid.New().String(),
				}
				_, err := app.CreateEvent(context.Background(), event)
				require.NoError(t, err)
				return fmt.Sprintf("/v1/events/daily?date=%v", url.QueryEscape(time.Now().Format(time.RFC3339)))
			},
			Method: http.MethodGet,
			Action: handler.GetDailyEvents,
			CheckResult: func(t *testing.T, res *httptest.ResponseRecorder, req *http.Request, st storage.Storage) { //nolint:thelper,lll
				var data EventsResponse
				err := json.Unmarshal(res.Body.Bytes(), &data)
				require.NoError(t, err)
				require.Equal(t, 1, len(data.Events), fmt.Sprintf("expected 1, found: %v", data.Events))
			},
		},
		"check GetWeekly": {
			URL: func(t *testing.T, app *app.App) string { //nolint:thelper
				event := storage.Event{
					Title:   "GetWeekly",
					Date:    time.Now().AddDate(0, 0, 7),
					EndDate: time.Now().AddDate(0, 0, 9),
					UserID:  uuid.New().String(),
				}
				_, err := app.CreateEvent(context.Background(), event)
				require.NoError(t, err)
				return fmt.Sprintf("/v1/events/weekly?date=%v", url.QueryEscape(time.Now().AddDate(0, 0, 3).Format(time.RFC3339)))
			},
			Method: http.MethodGet,
			Action: handler.GetDailyEvents,
			CheckResult: func(t *testing.T, res *httptest.ResponseRecorder, req *http.Request, st storage.Storage) { //nolint:thelper,lll
				var data EventsResponse
				err := json.Unmarshal(res.Body.Bytes(), &data)
				require.NoError(t, err)
				require.Equal(t, 1, len(data.Events), fmt.Sprintf("expected 1, found: %v", data.Events))
			},
		},
		"check GetMonthly": {
			URL: func(t *testing.T, app *app.App) string { //nolint:thelper
				event := storage.Event{
					Title:   "GetMonthly",
					Date:    time.Now().AddDate(0, 1, 0),
					EndDate: time.Now().AddDate(0, 1, 1),
					UserID:  uuid.New().String(),
				}
				_, err := app.CreateEvent(context.Background(), event)
				require.NoError(t, err)
				return "/v1/events/monthly"
			},
			Body:   fmt.Sprintf(`{"date": "%v"}`, time.Now().AddDate(0, 0, 20).Format(time.RFC3339)),
			Method: http.MethodGet,
			Action: handler.GetDailyEvents,
			CheckResult: func(t *testing.T, res *httptest.ResponseRecorder, req *http.Request, st storage.Storage) { //nolint:thelper,lll
				var data EventsResponse
				err := json.Unmarshal(res.Body.Bytes(), &data)
				require.NoError(t, err)
				require.Equal(t, 1, len(data.Events), fmt.Sprintf("expected 1, found: %v", data.Events))
			},
		},
	}

	for caseName, tc := range testcases {
		t.Run(caseName, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(tc.Action)) //nolint:unconvert
			defer testServer.Close()
			reader := strings.NewReader(tc.Body)
			req := httptest.NewRequest(tc.Method, tc.URL(t, application), reader)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			tc.CheckResult(t, res, req, st)
		})
	}
}

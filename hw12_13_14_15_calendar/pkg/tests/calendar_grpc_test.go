package integration_test

import (
	"context"
	"fmt"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	v1grpc "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api/v1/grpc"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	internalsql "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/pkg/rmq"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

var _ = Describe("Calendar GRPC", Ordered, func() {
	var currentEvent, weekAgoEvent, monthAgoEvent, yearAgoEvent *v1grpc.EventRequestValue
	var currentEventRes *v1grpc.EventIDResponse
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -8)
	monthAgo := now.AddDate(0, -1, 1)
	yearAgo := now.AddDate(-1, 0, -2)
	halfYearAgo := now.AddDate(0, -6, 1)

	conn, err := grpc.DialContext(ctx, grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(GinkgoT(), err)
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
	st := internalsql.New(&config)
	require.NoError(GinkgoT(), st.Connect(ctx))
	field := reflect.Indirect(reflect.ValueOf(st)).FieldByName("db")
	db := getUnexportedField(field).(*sqlx.DB)

	client := v1grpc.NewEventServiceClient(conn)

	BeforeAll(func() {
		_, err := db.Exec(`TRUNCATE TABLE events`)
		require.NoError(GinkgoT(), err)

		currentEvent = &v1grpc.EventRequestValue{
			Title:      gofakeit.Hobby(),
			Date:       timestamppb.New(now),
			EndDate:    timestamppb.New(now.Add(2 * time.Hour)),
			Desc:       gofakeit.Phrase(),
			UserId:     gofakeit.UUID(),
			RemindDate: timestamppb.New(now.Add(-1 * time.Minute)),
		}

		weekAgoEvent = &v1grpc.EventRequestValue{
			Title:      gofakeit.Hobby(),
			Date:       timestamppb.New(weekAgo),
			EndDate:    timestamppb.New(weekAgo.Add(2 * time.Hour)),
			Desc:       gofakeit.Phrase(),
			UserId:     gofakeit.UUID(),
			RemindDate: timestamppb.New(weekAgo.Add(-1 * time.Hour)),
		}

		monthAgoEvent = &v1grpc.EventRequestValue{
			Title:      gofakeit.Hobby(),
			Date:       timestamppb.New(monthAgo),
			EndDate:    timestamppb.New(monthAgo.Add(2 * time.Hour)),
			Desc:       gofakeit.Phrase(),
			UserId:     gofakeit.UUID(),
			RemindDate: timestamppb.New(monthAgo.Add(-1 * time.Hour)),
		}

		yearAgoEvent = &v1grpc.EventRequestValue{
			Title:      gofakeit.Hobby(),
			Date:       timestamppb.New(yearAgo),
			EndDate:    timestamppb.New(yearAgo.Add(2 * time.Hour)),
			Desc:       gofakeit.Phrase(),
			UserId:     gofakeit.UUID(),
			RemindDate: timestamppb.New(yearAgo.Add(-1 * time.Hour)),
		}
	})

	Describe("CreateEvent", func() {
		It("add year ago event", func() {
			_, err := client.CreateEvent(ctx, yearAgoEvent)
			require.NoError(GinkgoT(), err)
		})
		It("add now event", func() {
			var err error
			currentEventRes, err = client.CreateEvent(ctx, currentEvent)
			require.NoError(GinkgoT(), err)
		})
		It("add week ago event", func() {
			_, err := client.CreateEvent(ctx, weekAgoEvent)
			require.NoError(GinkgoT(), err)
		})
		It("add month ago event", func() {
			_, err := client.CreateEvent(ctx, monthAgoEvent)
			require.NoError(GinkgoT(), err)
		})
	})

	Describe("Complex", func() {
		It("success result", func() {
			mq := rmq.MessageQueue{}
			err := mq.Connect(amqpAddr)
			require.NoError(GinkgoT(), err)
			defer mq.Close()

			_, err = client.UpdateEvent(ctx, &v1grpc.UpdateRequest{
				Id: currentEventRes.Id,
				Event: &v1grpc.EventRequestValue{
					Date:       timestamppb.New(now.Add(3 * time.Hour)),
					EndDate:    timestamppb.New(now.AddDate(0, 0, 1)),
					RemindDate: timestamppb.New(now.Add(-1 * time.Hour)),
				},
			})
			require.NoError(GinkgoT(), err)

			msgs, err := mq.Channel.Consume(
				"status",
				"calendar",
				true,
				false,
				false,
				false,
				nil,
			)
			require.NoError(GinkgoT(), err)

			ticker := time.NewTicker(2 * sleepDuration)
			defer ticker.Stop()
			results := make(chan string)
			go func() {
				for {
					select {
					case msg, ok := <-msgs:
						if !ok {
							return
						}
						results <- string(msg.Body)
						return
					case <-ctx.Done():
						return
					}
				}
			}()

			select {
			case <-ticker.C:
				cancel()
			case data := <-results:
				require.Contains(GinkgoT(), data, "Successful send")
			case <-ctx.Done():
			}

			events, err := client.GetDailyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(yearAgo),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), 0, len(events.GetEvents()), fmt.Sprintf("Expected 0, found: %v", events.GetEvents()))
		})
	})

	Describe("GetWeeklyEvents", func() {
		It("success result", func() {
			events, err := client.GetWeeklyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(weekAgo),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "", events.Error)
			require.Equal(GinkgoT(), 1, len(events.GetEvents()))
		})
		It("empty result", func() {
			events, err := client.GetWeeklyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(halfYearAgo),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "", events.Error)
			require.Equal(GinkgoT(), 0, len(events.GetEvents()))
		})
	})

	Describe("GetMonthlyEvents", func() {
		It("success result", func() {
			events, err := client.GetMonthlyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(monthAgo),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "", events.Error)
			require.Equal(GinkgoT(), 3, len(events.GetEvents()))
		})
		It("empty result", func() {
			events, err := client.GetMonthlyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(halfYearAgo),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "", events.Error)
			require.Equal(GinkgoT(), 0, len(events.GetEvents()))
		})
	})
})

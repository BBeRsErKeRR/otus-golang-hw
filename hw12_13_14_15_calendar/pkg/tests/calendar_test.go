package integration_test

import (
	"context"
	"time"

	v1grpc "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/api/v1/grpc"
	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("Calendar GRPC", func() {
	var currentEvent, weekAgoEvent, monthAgoEvent, yearAgoEvent *v1grpc.EventRequestValue
	var currentEventRes *v1grpc.EventIDResponse
	ctx := context.Background()
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, -1, 1)
	yearAgo := now.AddDate(-1, 0, 1)
	halfYearAgo := now.AddDate(0, -6, 1)

	conn, err := grpc.DialContext(ctx, "0.0.0.0:5080",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(GinkgoT(), err)
	client := v1grpc.NewEventServiceClient(conn)

	BeforeEach(func() {
		currentEvent = &v1grpc.EventRequestValue{
			Title:      gofakeit.Hobby(),
			Date:       timestamppb.New(now),
			EndDate:    timestamppb.New(now.Add(2 * time.Hour)),
			Desc:       gofakeit.Phrase(),
			UserId:     gofakeit.UUID(),
			RemindDate: timestamppb.New(now.Add(1 * time.Hour)),
		}

		weekAgoEvent = &v1grpc.EventRequestValue{
			Title:      gofakeit.Hobby(),
			Date:       timestamppb.New(weekAgo),
			EndDate:    timestamppb.New(weekAgo.Add(2 * time.Hour)),
			Desc:       gofakeit.Phrase(),
			UserId:     gofakeit.UUID(),
			RemindDate: timestamppb.New(weekAgo.Add(1 * time.Hour)),
		}

		monthAgoEvent = &v1grpc.EventRequestValue{
			Title:      gofakeit.Hobby(),
			Date:       timestamppb.New(monthAgo),
			EndDate:    timestamppb.New(monthAgo.Add(2 * time.Hour)),
			Desc:       gofakeit.Phrase(),
			UserId:     gofakeit.UUID(),
			RemindDate: timestamppb.New(monthAgo.Add(1 * time.Hour)),
		}

		yearAgoEvent = &v1grpc.EventRequestValue{
			Title:      gofakeit.Hobby(),
			Date:       timestamppb.New(yearAgo),
			EndDate:    timestamppb.New(yearAgo.Add(2 * time.Hour)),
			Desc:       gofakeit.Phrase(),
			UserId:     gofakeit.UUID(),
			RemindDate: timestamppb.New(yearAgo.Add(1 * time.Hour)),
		}
	})

	Describe("CreateEvent", func() {
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
		It("add year ago event", func() {
			_, err := client.CreateEvent(ctx, yearAgoEvent)
			require.NoError(GinkgoT(), err)
		})
	})

	Describe("GetDailyEvents", func() {
		It("success result", func() {
			events, err := client.GetDailyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(now),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "", events.Error)
			require.Equal(GinkgoT(), 1, len(events.GetEvents()))
		})
		It("empty result", func() {
			events, err := client.GetDailyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(halfYearAgo),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "", events.Error)
			require.Equal(GinkgoT(), 0, len(events.GetEvents()))
		})
	})

	Describe("GetWeeklyEvents", func() {
		It("success result", func() {
			events, err := client.GetWeeklyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(weekAgo),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), "", events.Error)
			require.Equal(GinkgoT(), 2, len(events.GetEvents()))
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

	Describe("UpdateEvent", func() {
		It("success result", func() {
			_, err := client.UpdateEvent(ctx, &v1grpc.UpdateRequest{
				Id: currentEventRes.Id,
				Event: &v1grpc.EventRequestValue{
					Title:      gofakeit.Fruit(),
					Date:       timestamppb.New(now.AddDate(0, 0, 5)),
					EndDate:    timestamppb.New(now.AddDate(0, 0, 6)),
					RemindDate: timestamppb.New(now.AddDate(0, 0, 5).Add(2 * time.Hour)),
				},
			})
			require.NoError(GinkgoT(), err)

			events, err := client.GetDailyEvents(ctx, &v1grpc.EventsRequest{
				Date: timestamppb.New(now),
			})
			require.NoError(GinkgoT(), err)
			require.Equal(GinkgoT(), 0, len(events.GetEvents()))
		})
	})
})

package v1grpc

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestV1GRPCHandlers(t *testing.T) {
	ctx := context.Background()
	loggerConfig := logger.Config{
		Level:    "debug",
		OutPaths: []string{},
		ErrPaths: []string{},
	}
	logger, err := logger.New(&loggerConfig)
	require.NoError(t, err)
	st := memorystorage.New()
	application := app.New(logger, *storage.New(st))

	lis := bufconn.Listen(101024 * 1024)
	require.NoError(t, err)
	baseServer := grpc.NewServer()
	RegisterEventServiceServer(baseServer, NewHandler(application, logger))
	defer baseServer.Stop()
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			require.NoError(t, err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := NewEventServiceClient(conn)

	t.Run("crud case", func(t *testing.T) {
		crResp, err := client.CreateEvent(ctx,
			&Event{
				Title:   "Created",
				Date:    timestamppb.Now(),
				EndDate: timestamppb.New(time.Now().AddDate(0, 0, 1)),
				UserID:  uuid.New().String(),
			})
		require.NoError(t, err)
		updResp, err := client.UpdateEvent(ctx,
			&UpdateRequest{
				Id: crResp.Id,
				Event: &Event{
					Title:   "UPDATED",
					Date:    timestamppb.Now(),
					EndDate: timestamppb.New(time.Now().AddDate(0, 1, 0)),
					UserID:  uuid.New().String(),
				},
			},
		)
		require.NoError(t, err)
		require.Equal(t, "success", updResp.Msg)
		event, err := st.GetEvent(ctx, crResp.Id)
		require.NoError(t, err)
		require.Equal(t, "UPDATED", event.Title)
		delResp, err := client.DeleteEvent(ctx,
			&EventID{
				Id: crResp.Id,
			},
		)
		require.NoError(t, err)
		require.Equal(t, "success", delResp.Msg)
		_, err = st.GetEvent(ctx, crResp.Id)
		require.ErrorIs(t, storage.ErrNotExist, err)
	})
}

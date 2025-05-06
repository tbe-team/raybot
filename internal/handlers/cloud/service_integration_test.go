package cloud_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/handlers/cloud"
	"github.com/tbe-team/raybot/internal/handlers/cloud/cloudtest"
	"github.com/tbe-team/raybot/internal/logging"
	commandmocks "github.com/tbe-team/raybot/internal/services/command/mocks"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

func TestIntegrationService_ConnectAndDisconnect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	port, stop := cloudtest.SetupTestCloudServer(t)
	defer stop()

	cfg := config.Cloud{
		Address: fmt.Sprintf("localhost:%d", port),
	}
	log := logging.NewNoopLogger()
	publisher := newFakePublisher()
	commandService := commandmocks.NewFakeService(t)
	service := cloud.New(
		cfg,
		log,
		publisher,
		commandService,
		cloud.WithConnectTimeout(50*time.Millisecond),
	)

	cleanup, err := service.Run(context.Background())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cleanup())
	}()

	// Wait for the cloud to connect
	msg := publisher.WaitForPublish(events.CloudConnectedTopic, 2*time.Second)
	require.NotNil(t, msg)

	// stop the tunnel server
	stop()

	// Wait for the cloud to disconnect
	msg = publisher.WaitForPublish(events.CloudDisconnectedTopic, 2*time.Second)
	require.NotNil(t, msg)
}

func TestIntegrationService_CloudServerNotRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := config.Cloud{
		Address: "localhost:1234",
	}
	log := logging.NewNoopLogger()
	publisher := newFakePublisher()
	commandService := commandmocks.NewFakeService(t)
	service := cloud.New(
		cfg,
		log,
		publisher,
		commandService,
		cloud.WithConnectTimeout(50*time.Millisecond),
	)

	cleanup, err := service.Run(context.Background())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cleanup())
	}()

	// Wait for the disconnected event
	msg := publisher.WaitForPublish(events.CloudDisconnectedTopic, 2*time.Second)
	require.NotNil(t, msg)
}

type fakePublisher struct {
	events chan *eventbus.Message
}

func newFakePublisher() *fakePublisher {
	return &fakePublisher{
		events: make(chan *eventbus.Message, 1),
	}
}

func (f fakePublisher) Publish(_ string, msg *eventbus.Message) {
	select {
	case f.events <- msg:
	default:
	}
}

func (f fakePublisher) WaitForPublish(_ string, timeout time.Duration) *eventbus.Message {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case msg := <-f.events:
		return msg
	case <-timer.C:
		return nil
	}
}

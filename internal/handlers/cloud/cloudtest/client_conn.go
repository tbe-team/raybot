package cloudtest

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jhump/grpctunnel"
	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/handlers/cloud"
	"github.com/tbe-team/raybot/internal/logging"
	"github.com/tbe-team/raybot/internal/services/appstate/appstateimpl"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/command/commandimpl"
	"github.com/tbe-team/raybot/internal/services/command/executor"
	"github.com/tbe-team/raybot/internal/services/command/processinglockimpl"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/validator"
)

type TunnelTestEnv struct {
	TunnelChannel grpctunnel.TunnelChannel

	CommandService command.Service
}

// SetupTunnelTestEnv sets up a temporary cloud server with gRPC reverse tunnel,
// initializes the required services, and waits for the tunnel to be established.
// It returns a TunnelTestEnv that contains all things needed for integration tests,
// so test clients can use it to send gRPC calls through the tunnel.
//
// This is used for integration testing where the client connects to the cloud through
// a reverse tunnel.
//
// Automatically cleans up resources after the test.
func SetupTunnelTestEnv(t *testing.T) TunnelTestEnv {
	t.Helper()

	var tc grpctunnel.TunnelChannel
	wg := sync.WaitGroup{}

	wg.Add(1)
	port, stop := SetupTestCloudServer(
		t,
		WithTunnelServiceHandlerOpts(grpctunnel.TunnelServiceHandlerOptions{
			OnReverseTunnelOpen: func(tunnelChannel grpctunnel.TunnelChannel) {
				tc = tunnelChannel
				wg.Done()
			},
		}),
	)

	db, err := db.NewTestDB()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate())
	queries := sqlc.New()
	log := logging.NewNoopLogger()
	bus := eventbus.NewNoopEventBus()
	validator := validator.New()
	commandService := commandimpl.NewService(
		config.DeleteOldCommand{},
		log,
		validator,
		bus,
		commandimpl.NewCommandRepository(db, queries),
		appstateimpl.NewAppStateRepository(),
		processinglockimpl.New(),
		executor.NewNoopRouter(),
	)

	cloudSvc := cloud.New(
		config.Cloud{
			Address: fmt.Sprintf("localhost:%d", port),
		},
		log,
		bus,
		commandService,
		cloud.WithConnectTimeout(500*time.Millisecond),
	)
	cleanupCloudSvc, err := cloudSvc.Run(context.Background())
	require.NoError(t, err)

	// Wait for the tunnel to be opened
	wg.Wait()

	t.Cleanup(func() {
		require.NoError(t, cleanupCloudSvc())
		stop()
		require.NoError(t, db.Close())
	})

	return TunnelTestEnv{
		CommandService: commandService,
		TunnelChannel:  tc,
	}
}

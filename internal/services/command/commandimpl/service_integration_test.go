package commandimpl

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/logging"
	"github.com/tbe-team/raybot/internal/services/appstate/appstateimpl"
	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/services/command/processinglockimpl"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
	"github.com/tbe-team/raybot/pkg/eventbus"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/validator"
)

func TestIntegrationCommandService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run(`Created 3 commands from cloud, 1 will be in PROCESSING, 2 will be in QUEUED,
		CancelActiveCloudCommands should cancel 3 commands and there is no running command`,
		func(t *testing.T) {
			log := logging.NewNoopLogger()
			db, err := db.NewTestDB()
			require.NoError(t, err)
			defer func() {
				require.NoError(t, db.Close())
			}()
			err = db.AutoMigrate()
			require.NoError(t, err)
			queries := sqlc.New()
			runningCmdRepository := newRunningCmdRepository()
			commandRepository := NewCommandRepository(db, queries)
			router := newBlockingExecutorRouter(commandRepository)

			commandService := Service{
				deleteOldCmdCfg:      config.DeleteOldCommand{},
				log:                  log,
				validator:            validator.New(),
				publisher:            eventbus.NewInProcEventBus(log),
				runningCmdRepository: runningCmdRepository,
				commandRepository:    commandRepository,
				appStateRepository:   appstateimpl.NewAppStateRepository(),
				processingLock:       processinglockimpl.New(),
				executorRouter:       router,
			}

			// Create 2 commands from cloud
			cmd1, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
				Source: command.SourceCloud,
				Inputs: command.MoveForwardInputs{},
			})
			require.NoError(t, err)

			cmd2, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
				Source: command.SourceCloud,
				Inputs: command.MoveForwardInputs{},
			})
			require.NoError(t, err)

			cmd3, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
				Source: command.SourceCloud,
				Inputs: command.MoveForwardInputs{},
			})
			require.NoError(t, err)

			// Because we don't handle event in service layer
			// so we simulate handle [events.CommandCreatedEvent] here
			go func() {
				// This should block until executor router unblock
				if err := commandService.ExecuteCreatedCommand(context.Background(), command.ExecuteCreatedCommandParams{
					CommandID: cmd1.ID,
				}); err != nil {
					t.Errorf("failed to execute command: %v", err)
				}
			}()

			// Block until we have a running command (cause by ExecuteCreatedCommand)
			require.Eventually(t, func() bool {
				return runningCmdRepository.Get() != nil
			}, 500*time.Millisecond, 10*time.Millisecond)

			// Get command to validate state
			cmd1, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd1.ID,
			})
			require.NoError(t, err)
			require.Equal(t, command.StatusProcessing, cmd1.Status)

			cmd2, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd2.ID,
			})
			require.NoError(t, err)
			require.Equal(t, command.StatusQueued, cmd2.Status)

			cmd3, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd3.ID,
			})
			require.NoError(t, err)
			require.Equal(t, command.StatusQueued, cmd3.Status)

			// Current running should be cmd1
			runningCmd := runningCmdRepository.Get()
			require.NotNil(t, runningCmd)
			require.Equal(t, cmd1.ID, runningCmd.ID)

			err = commandService.CancelActiveCloudCommands(context.Background())
			require.NoError(t, err)

			// Block until no running command (cause by CancelActiveCloudCommands)
			require.Eventually(t, func() bool {
				return runningCmdRepository.Get() == nil
			}, 500*time.Millisecond, 10*time.Millisecond)

			// Get command to validate state
			cmd1, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd1.ID,
			})
			require.NoError(t, err)
			require.Equal(t, command.StatusCanceled, cmd1.Status)

			cmd2, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd2.ID,
			})
			require.NoError(t, err)
			require.Equal(t, command.StatusCanceled, cmd2.Status)

			cmd3, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd3.ID,
			})
			require.NoError(t, err)
			require.Equal(t, command.StatusCanceled, cmd3.Status)
		},
	)

	t.Run(`Created 2 commands, 1 will be in PROCESSING, 1 will be in QUEUED,
		CancelCurrentProcessingCommand should cancel the running command and the next command should be in PROCESSING`,
		func(t *testing.T) {
			log := logging.NewNoopLogger()
			db, err := db.NewTestDB()
			require.NoError(t, err)
			defer func() {
				require.NoError(t, db.Close())
			}()
			err = db.AutoMigrate()
			require.NoError(t, err)
			queries := sqlc.New()
			commandRepository := NewCommandRepository(db, queries)
			runningCmdRepository := newRunningCmdRepository()
			commandService := Service{
				deleteOldCmdCfg:      config.DeleteOldCommand{},
				log:                  log,
				validator:            validator.New(),
				publisher:            eventbus.NewInProcEventBus(log),
				runningCmdRepository: runningCmdRepository,
				commandRepository:    commandRepository,
				appStateRepository:   appstateimpl.NewAppStateRepository(),
				processingLock:       processinglockimpl.New(),
				executorRouter:       newBlockingExecutorRouter(commandRepository),
			}

			// Create 2 commands
			cmd1, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
				Source: command.SourceCloud,
				Inputs: command.MoveForwardInputs{},
			})
			require.NoError(t, err)

			cmd2, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
				Source: command.SourceCloud,
				Inputs: command.MoveForwardInputs{},
			})
			require.NoError(t, err)

			// Because we don't handle event in service layer
			// so we simulate handle [events.CommandCreatedEvent] here
			go func() {
				// This should block until executor router unblock
				if err := commandService.ExecuteCreatedCommand(context.Background(), command.ExecuteCreatedCommandParams{
					CommandID: cmd1.ID,
				}); err != nil {
					t.Errorf("failed to execute command: %v", err)
				}
			}()

			// Block until we have a running command (cause by ExecuteCreatedCommand)
			require.Eventually(t, func() bool {
				return runningCmdRepository.Get() != nil
			}, 500*time.Millisecond, 10*time.Millisecond)

			// Get command to validate state
			cmd1, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd1.ID,
			})
			require.NoError(t, err)
			require.Equal(t, command.StatusProcessing, cmd1.Status)

			cmd2, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd2.ID,
			})
			require.NoError(t, err)
			require.Equal(t, command.StatusQueued, cmd2.Status)

			err = commandService.CancelCurrentProcessingCommand(context.Background())
			require.NoError(t, err)

			// Block until next command is in running (cause by ExecuteCreatedCommand)
			require.Eventually(t, func() bool {
				cmd := runningCmdRepository.Get()
				return cmd != nil && cmd.ID == cmd2.ID && cmd.Status == command.StatusProcessing
			}, 500*time.Millisecond, 10*time.Millisecond)
		})

	t.Run(`When create CommandService it should run cancel QUEUED and PROCESSING commands`, func(t *testing.T) {
		log := logging.NewNoopLogger()
		db, err := db.NewTestDB()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, db.Close())
		}()
		err = db.AutoMigrate()
		require.NoError(t, err)
		queries := sqlc.New()
		commandRepository := NewCommandRepository(db, queries)

		// Create 2 commands 1 in QUEUED and 1 in PROCESSING
		cmd1, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status: command.StatusQueued,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		cmd2, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status: command.StatusProcessing,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		commandService := NewService(
			config.DeleteOldCommand{},
			log,
			validator.New(),
			eventbus.NewInProcEventBus(log),
			commandRepository,
			appstateimpl.NewAppStateRepository(),
			processinglockimpl.New(),
			newBlockingExecutorRouter(commandRepository),
		)

		require.Eventually(t, func() bool {
			cmd1, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
				CommandID: cmd1.ID,
			})
			require.NoError(t, err)
			return cmd1.Status == command.StatusCanceled
		}, 500*time.Millisecond, 10*time.Millisecond)

		cmd2, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
			CommandID: cmd2.ID,
		})
		require.NoError(t, err)
		require.Equal(t, command.StatusCanceled, cmd2.Status)
	})

	t.Run(`Get current processing command should return the command in PROCESSING status`, func(t *testing.T) {
		log := logging.NewNoopLogger()
		db, err := db.NewTestDB()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, db.Close())
		}()
		err = db.AutoMigrate()
		require.NoError(t, err)
		queries := sqlc.New()
		commandRepository := NewCommandRepository(db, queries)
		runningCmdRepository := newRunningCmdRepository()
		commandService := Service{
			deleteOldCmdCfg:      config.DeleteOldCommand{},
			log:                  log,
			validator:            validator.New(),
			publisher:            eventbus.NewInProcEventBus(log),
			runningCmdRepository: runningCmdRepository,
			commandRepository:    commandRepository,
			appStateRepository:   appstateimpl.NewAppStateRepository(),
			processingLock:       processinglockimpl.New(),
			executorRouter:       newBlockingExecutorRouter(commandRepository),
		}

		cmd1, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status: command.StatusProcessing,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		cmd, err := commandService.GetCurrentProcessingCommand(context.Background())
		require.NoError(t, err)
		require.Equal(t, cmd1.ID, cmd.ID)
		require.Equal(t, command.StatusProcessing, cmd.Status)
	})

	t.Run(`List commands should return list of commands`, func(t *testing.T) {
		log := logging.NewNoopLogger()
		db, err := db.NewTestDB()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, db.Close())
		}()
		err = db.AutoMigrate()
		require.NoError(t, err)
		queries := sqlc.New()
		commandRepository := NewCommandRepository(db, queries)
		runningCmdRepository := newRunningCmdRepository()
		commandService := Service{
			deleteOldCmdCfg:      config.DeleteOldCommand{},
			log:                  log,
			validator:            validator.New(),
			publisher:            eventbus.NewInProcEventBus(log),
			runningCmdRepository: runningCmdRepository,
			commandRepository:    commandRepository,
			appStateRepository:   appstateimpl.NewAppStateRepository(),
			processingLock:       processinglockimpl.New(),
			executorRouter:       newBlockingExecutorRouter(commandRepository),
		}

		cmd1, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status: command.StatusProcessing,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		cmd2, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status: command.StatusQueued,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		commands, err := commandService.ListCommands(context.Background(), command.ListCommandsParams{
			PagingParams: paging.NewParams(paging.Page(1), paging.PageSize(10)),
		})
		require.NoError(t, err)

		var ids []int64
		for _, cmd := range commands.Items {
			ids = append(ids, cmd.ID)
		}
		require.Contains(t, ids, cmd1.ID)
		require.Contains(t, ids, cmd2.ID)
	})

	t.Run(`Lock processing command successfully`, func(t *testing.T) {
		/*
		   Given a queue with three commands: cmd1, cmd2, and cmd3, created in that order.
		   - cmd1 enters the PROCESSING state (assumed to be blocking).
		   - cmd2 and cmd3 remain in the QUEUED state.

		   When a call is made to forcibly lock the processing command:
		   - cmd1 should be cancelled.
		   - cmd2 and cmd3 should remain in the QUEUED state.
		*/

		log := logging.NewNoopLogger()
		db, err := db.NewTestDB()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, db.Close())
		}()
		err = db.AutoMigrate()
		require.NoError(t, err)
		queries := sqlc.New()
		commandRepository := NewCommandRepository(db, queries)
		runningCmdRepository := newRunningCmdRepository()
		commandService := Service{
			deleteOldCmdCfg:      config.DeleteOldCommand{},
			log:                  log,
			validator:            validator.New(),
			publisher:            eventbus.NewInProcEventBus(log),
			runningCmdRepository: runningCmdRepository,
			commandRepository:    commandRepository,
			appStateRepository:   appstateimpl.NewAppStateRepository(),
			processingLock:       processinglockimpl.New(),
			executorRouter:       newBlockingExecutorRouter(commandRepository),
		}

		cmd1, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
			Source: command.SourceCloud,
			Inputs: command.MoveForwardInputs{},
		})
		require.NoError(t, err)

		cmd2, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
			Source: command.SourceCloud,
			Inputs: command.MoveForwardInputs{},
		})
		require.NoError(t, err)

		cmd3, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
			Source: command.SourceCloud,
			Inputs: command.MoveForwardInputs{},
		})
		require.NoError(t, err)

		// Because we don't handle event in service layer
		// so we simulate handle [events.CommandCreatedEvent] here
		go func() {
			// This should block until executor router unblock
			if err := commandService.ExecuteCreatedCommand(context.Background(), command.ExecuteCreatedCommandParams{
				CommandID: cmd1.ID,
			}); err != nil {
				t.Errorf("failed to execute command: %v", err)
			}
		}()

		// Block until we have a running command (cause by ExecuteCreatedCommand)
		require.Eventually(t, func() bool {
			return runningCmdRepository.Get() != nil
		}, 500*time.Millisecond, 10*time.Millisecond)

		err = commandService.LockProcessingCommand(context.Background())
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			return runningCmdRepository.Get() == nil
		}, 500*time.Millisecond, 10*time.Millisecond)

		cmd1, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
			CommandID: cmd1.ID,
		})
		require.NoError(t, err)
		require.Equal(t, command.StatusCanceled, cmd1.Status)

		cmd2, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
			CommandID: cmd2.ID,
		})
		require.NoError(t, err)
		require.Equal(t, command.StatusQueued, cmd2.Status)

		cmd3, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
			CommandID: cmd3.ID,
		})
		require.NoError(t, err)
		require.Equal(t, command.StatusQueued, cmd3.Status)
	})

	t.Run(`Execute created command should block until processing lock is released`, func(t *testing.T) {
		/*
		   Given a cmd1 and processing lock is acquired

		   When ExecuteCreatedCommand is called with cmd1.ID
		   - It should block until processing lock is released
		   - cmd1 status should be QUEUED before processing lock is released
		*/

		log := logging.NewNoopLogger()
		db, err := db.NewTestDB()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, db.Close())
		}()
		err = db.AutoMigrate()
		require.NoError(t, err)
		queries := sqlc.New()
		commandRepository := NewCommandRepository(db, queries)
		runningCmdRepository := newRunningCmdRepository()
		commandService := Service{
			deleteOldCmdCfg:      config.DeleteOldCommand{},
			log:                  log,
			validator:            validator.New(),
			publisher:            eventbus.NewInProcEventBus(log),
			runningCmdRepository: runningCmdRepository,
			commandRepository:    commandRepository,
			appStateRepository:   appstateimpl.NewAppStateRepository(),
			processingLock:       processinglockimpl.New(),
			executorRouter:       newBlockingExecutorRouter(commandRepository),
		}

		cmd1, err := commandService.CreateCommand(context.Background(), command.CreateCommandParams{
			Source: command.SourceCloud,
			Inputs: command.MoveForwardInputs{},
		})
		require.NoError(t, err)

		err = commandService.LockProcessingCommand(context.Background())
		require.NoError(t, err)

		errCh := make(chan error)
		go func() {
			err := commandService.ExecuteCreatedCommand(context.Background(), command.ExecuteCreatedCommandParams{
				CommandID: cmd1.ID,
			})
			errCh <- err
		}()

		time.Sleep(100 * time.Millisecond)

		err = commandService.UnlockProcessingCommand(context.Background())
		require.NoError(t, err)

		cmd1, err = commandService.GetCommandByID(context.Background(), command.GetCommandByIDParams{
			CommandID: cmd1.ID,
		})
		require.NoError(t, err)
		require.Equal(t, command.StatusQueued, cmd1.Status)
	})
}

// blockingExecutorRouter is a mock Router used for testing.
// It simulates a blocked executor and updates the command status to Canceled
// if the context is canceled. Useful for testing cancelation behavior in the command service.
type blockingExecutorRouter struct {
	doneCh            chan struct{}
	commandRepository command.Repository
}

func newBlockingExecutorRouter(commandRepository command.Repository) *blockingExecutorRouter {
	return &blockingExecutorRouter{
		doneCh:            make(chan struct{}),
		commandRepository: commandRepository,
	}
}

func (r *blockingExecutorRouter) Route(ctx context.Context, cmd command.Command) error {
	select {
	case <-r.doneCh:
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			_, err := r.commandRepository.UpdateCommand(context.Background(), command.UpdateCommandParams{
				ID:        cmd.ID,
				Status:    command.StatusCanceled,
				SetStatus: true,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *blockingExecutorRouter) UnBlock() {
	close(r.doneCh)
}

package commandimpl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/logging"
	"github.com/tbe-team/raybot/internal/services/command"
	commandmocks "github.com/tbe-team/raybot/internal/services/command/mocks"
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

	t.Run(`When create CommandService it should cancel all QUEUED and PROCESSING commands`, func(t *testing.T) {
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

		runningCmdRepository := NewRunningCmdRepository()
		commandService := NewService(
			config.DeleteOldCommand{},
			log,
			validator.New(),
			eventbus.NewInProcEventBus(log),
			runningCmdRepository,
			commandRepository,
			processinglockimpl.New(),
			commandmocks.NewFakeExecutorService(t),
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
		runningCmdRepository := NewRunningCmdRepository()
		commandService := Service{
			deleteOldCmdCfg:      config.DeleteOldCommand{},
			log:                  log,
			validator:            validator.New(),
			publisher:            eventbus.NewInProcEventBus(log),
			runningCmdRepository: runningCmdRepository,
			commandRepository:    commandRepository,
			processingLock:       processinglockimpl.New(),
			executorService:      commandmocks.NewFakeExecutorService(t),
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
		runningCmdRepository := NewRunningCmdRepository()
		commandService := Service{
			deleteOldCmdCfg:      config.DeleteOldCommand{},
			log:                  log,
			validator:            validator.New(),
			publisher:            eventbus.NewInProcEventBus(log),
			runningCmdRepository: runningCmdRepository,
			commandRepository:    commandRepository,
			processingLock:       processinglockimpl.New(),
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

	t.Run("Delete command by id should not delete command with status PROCESSING", func(t *testing.T) {
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
		commandService := Service{
			log:               log,
			validator:         validator.New(),
			publisher:         eventbus.NewInProcEventBus(log),
			commandRepository: commandRepository,
		}

		cmd, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status: command.StatusProcessing,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		err = commandService.DeleteCommandByID(context.Background(), command.DeleteCommandByIDParams{
			CommandID: cmd.ID,
		})
		require.ErrorIs(t, err, command.ErrCommandInProcessingCanNotBeDeleted)
	})

	t.Run("Delete old commands should delete commands follow by cutoff time and not in status QUEUED or PROCESSING", func(t *testing.T) {
		db, err := db.NewTestDB()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, db.Close())
		}()
		err = db.AutoMigrate()
		require.NoError(t, err)
		queries := sqlc.New()
		threshold := 1 * time.Hour
		commandRepository := NewCommandRepository(db, queries)
		commandService := Service{
			deleteOldCmdCfg: config.DeleteOldCommand{
				Threshold: threshold,
			},
			validator:         validator.New(),
			commandRepository: commandRepository,
		}

		cmd1, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status:    command.StatusSucceeded,
			Type:      command.CommandTypeStopMovement,
			CreatedAt: time.Now().Add(-2 * threshold),
		})
		require.NoError(t, err)

		cmd2, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status:    command.StatusProcessing,
			Type:      command.CommandTypeStopMovement,
			CreatedAt: time.Now(),
		})
		require.NoError(t, err)

		cmd3, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status:    command.StatusQueued,
			Type:      command.CommandTypeStopMovement,
			CreatedAt: time.Now().Add(-2 * threshold),
		})
		require.NoError(t, err)

		err = commandService.DeleteOldCommands(context.Background())
		require.NoError(t, err)

		cmd1, err = commandRepository.GetCommandByID(context.Background(), cmd1.ID)
		t.Log(err)
		require.ErrorIs(t, err, command.ErrCommandNotFound)

		cmd2, err = commandRepository.GetCommandByID(context.Background(), cmd2.ID)
		require.NoError(t, err)
		require.Equal(t, command.StatusProcessing, cmd2.Status)

		cmd3, err = commandRepository.GetCommandByID(context.Background(), cmd3.ID)
		require.NoError(t, err)
		require.Equal(t, command.StatusQueued, cmd3.Status)
	})

	t.Run("Cancel all running commands should cancel current running command and all queued commands", func(t *testing.T) {
		db, err := db.NewTestDB()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, db.Close())
		}()
		err = db.AutoMigrate()
		require.NoError(t, err)
		queries := sqlc.New()
		threshold := 1 * time.Hour
		runningCmdRepository := NewRunningCmdRepository()
		commandRepository := NewCommandRepository(db, queries)
		commandService := Service{
			deleteOldCmdCfg: config.DeleteOldCommand{
				Threshold: threshold,
			},
			validator:            validator.New(),
			commandRepository:    commandRepository,
			runningCmdRepository: runningCmdRepository,
			processingLock:       processinglockimpl.New(),
		}

		cmd1, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status: command.StatusProcessing,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		err = runningCmdRepository.Add(context.Background(), command.NewCancelableCommand(context.Background(), cmd1))
		require.NoError(t, err)

		cmd2, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Status: command.StatusQueued,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		err = commandService.CancelAllRunningCommands(context.Background())
		require.NoError(t, err)

		cmd1, err = commandRepository.GetCommandByID(context.Background(), cmd1.ID)
		require.NoError(t, err)
		require.Equal(t, command.StatusCanceled, cmd1.Status)

		runningCmd, err := runningCmdRepository.Get(context.Background())
		require.ErrorIs(t, err, command.ErrRunningCommandNotFound)
		require.Empty(t, runningCmd)

		cmd2, err = commandRepository.GetCommandByID(context.Background(), cmd2.ID)
		require.NoError(t, err)
		require.Equal(t, command.StatusCanceled, cmd2.Status)
	})

	t.Run("Cancel current running command successfully", func(t *testing.T) {
		runningCmdRepository := NewRunningCmdRepository()
		commandService := Service{
			runningCmdRepository: runningCmdRepository,
		}

		cmd := command.NewCommand(command.SourceApp, command.MoveForwardInputs{})
		err := runningCmdRepository.Add(context.Background(), command.NewCancelableCommand(context.Background(), cmd))
		require.NoError(t, err)

		err = commandService.CancelCurrentProcessingCommand(context.Background())
		require.NoError(t, err)

		runningCmd, err := runningCmdRepository.Get(context.Background())
		require.NoError(t, err)

		select {
		case <-runningCmd.Context().Done():
		case <-time.After(1 * time.Millisecond):
			require.Fail(t, "command should be canceled")
		}
	})

	t.Run("Cancel active cloud commands should cancel all QUEUED and PROCESSING commands created by the cloud", func(t *testing.T) {
		db, err := db.NewTestDB()
		require.NoError(t, err)
		defer func() {
			require.NoError(t, db.Close())
		}()
		err = db.AutoMigrate()
		require.NoError(t, err)
		queries := sqlc.New()
		runningCmdRepository := NewRunningCmdRepository()
		commandRepository := NewCommandRepository(db, queries)
		commandService := Service{
			validator:            validator.New(),
			commandRepository:    commandRepository,
			runningCmdRepository: runningCmdRepository,
			processingLock:       processinglockimpl.New(),
		}

		cmd1, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Source: command.SourceCloud,
			Status: command.StatusQueued,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		cmd2, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Source: command.SourceCloud,
			Status: command.StatusProcessing,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		err = runningCmdRepository.Add(context.Background(), command.NewCancelableCommand(context.Background(), cmd2))
		require.NoError(t, err)

		cmd3, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Source: command.SourceCloud,
			Status: command.StatusSucceeded,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		cmd4, err := commandRepository.CreateCommand(context.Background(), command.Command{
			Source: command.SourceApp,
			Status: command.StatusQueued,
			Type:   command.CommandTypeStopMovement,
		})
		require.NoError(t, err)

		err = commandService.CancelActiveCloudCommands(context.Background())
		require.NoError(t, err)

		cmd1, err = commandRepository.GetCommandByID(context.Background(), cmd1.ID)
		require.NoError(t, err)
		require.Equal(t, command.StatusCanceled, cmd1.Status)

		cmd2, err = commandRepository.GetCommandByID(context.Background(), cmd2.ID)
		require.NoError(t, err)
		require.Equal(t, command.StatusCanceled, cmd2.Status)

		cmd3, err = commandRepository.GetCommandByID(context.Background(), cmd3.ID)
		require.NoError(t, err)
		require.Equal(t, command.StatusSucceeded, cmd3.Status)

		cmd4, err = commandRepository.GetCommandByID(context.Background(), cmd4.ID)
		require.NoError(t, err)
		require.Equal(t, command.StatusQueued, cmd4.Status)
	})
}

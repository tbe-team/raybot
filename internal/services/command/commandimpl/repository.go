package commandimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/sync/errgroup"

	"github.com/tbe-team/raybot/internal/services/command"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/ptr"
)

type repository struct {
	db      db.DB
	queries *sqlc.Queries
}

func NewCommandRepository(db db.DB, queries *sqlc.Queries) command.Repository {
	return &repository{
		db:      db,
		queries: queries,
	}
}

func (r repository) ListCommands(ctx context.Context, params command.ListCommandsParams) (paging.List[command.Command], error) {
	query := sq.Select("*").
		From("commands").
		Limit(uint64(params.PagingParams.Limit())).
		Offset(uint64(params.PagingParams.Offset()))

	for _, s := range params.Sorts {
		query = s.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return paging.List[command.Command]{}, fmt.Errorf("failed to build query: %w", err)
	}

	countQuery := sq.Select("COUNT(*)").
		From("commands")

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return paging.List[command.Command]{}, fmt.Errorf("failed to build count query: %w", err)
	}

	var ret paging.List[command.Command]
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		rows, err := r.db.QueryContext(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf("query commands: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var row sqlc.Command
			if err := rows.Scan(
				&row.ID,
				&row.Type,
				&row.Status,
				&row.Source,
				&row.Inputs,
				&row.Error,
				&row.CompletedAt,
				&row.CreatedAt,
				&row.UpdatedAt,
			); err != nil {
				return fmt.Errorf("scan command: %w", err)
			}

			cmd, err := r.convertRowToCommand(row)
			if err != nil {
				return fmt.Errorf("convert row to command: %w", err)
			}
			ret.Items = append(ret.Items, cmd)
		}

		return nil
	})

	g.Go(func() error {
		if err := r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&ret.TotalItems); err != nil {
			return fmt.Errorf("scan count row: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return paging.List[command.Command]{}, fmt.Errorf("errgroup wait: %w", err)
	}

	return ret, nil
}

func (r repository) CreateCommand(ctx context.Context, commandArg command.Command) (command.Command, error) {
	inputsBytes, err := json.Marshal(commandArg.Inputs)
	if err != nil {
		return command.Command{}, fmt.Errorf("failed to marshal inputs: %w", err)
	}

	var completedAt *string
	if commandArg.CompletedAt != nil {
		completedAt = ptr.New(commandArg.CompletedAt.Format(time.RFC3339))
	}

	id, err := r.queries.CommandCreate(ctx, r.db, sqlc.CommandCreateParams{
		Type:        commandArg.Type.String(),
		Status:      commandArg.Status.String(),
		Source:      commandArg.Source.String(),
		Inputs:      string(inputsBytes),
		Error:       commandArg.Error,
		CreatedAt:   commandArg.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   commandArg.UpdatedAt.Format(time.RFC3339),
		CompletedAt: completedAt,
	})
	if err != nil {
		return command.Command{}, fmt.Errorf("queries create command: %w", err)
	}

	commandArg.ID = id

	return commandArg, nil
}

func (repository) convertRowToCommand(row sqlc.Command) (command.Command, error) {
	ret := command.Command{
		ID:     row.ID,
		Type:   command.CommandType(row.Type),
		Status: command.Status(row.Status),
		Source: command.Source(row.Source),
		Error:  row.Error,
	}
	var err error

	ret.Inputs, err = command.UnmarshalInputs(command.CommandType(row.Type), []byte(row.Inputs))
	if err != nil {
		return command.Command{}, fmt.Errorf("failed to unmarshal inputs: %w", err)
	}

	ret.CreatedAt, err = time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return command.Command{}, fmt.Errorf("failed to parse created at: %w", err)
	}

	ret.UpdatedAt, err = time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return command.Command{}, fmt.Errorf("failed to parse updated at: %w", err)
	}

	if row.CompletedAt != nil {
		completedAt, err := time.Parse(time.RFC3339, *row.CompletedAt)
		if err != nil {
			return command.Command{}, fmt.Errorf("failed to parse completed at: %w", err)
		}
		ret.CompletedAt = &completedAt
	}

	return ret, nil
}

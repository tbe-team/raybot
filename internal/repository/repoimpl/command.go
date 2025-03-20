package repoimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/sort"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var (
	ErrCommandNotFound = xerror.NotFound(nil, "command.notFound", "command not found")
)

type CommandRepository struct {
	queries *sqlc.Queries
}

func NewCommandRepository(queries *sqlc.Queries) *CommandRepository {
	return &CommandRepository{queries: queries}
}

func (r CommandRepository) ListCommands(ctx context.Context, db db.SQLDB, params paging.Params, sorts []sort.Sort) (paging.List[model.Command], error) {
	query := sq.Select("*").
		From("commands").
		Limit(uint64(params.Limit())).
		Offset(uint64(params.Offset()))

	for _, s := range sorts {
		query = s.Attach(query)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return paging.List[model.Command]{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return paging.List[model.Command]{}, fmt.Errorf("query commands: %w", err)
	}
	defer rows.Close()

	items := make([]model.Command, 0, params.Limit())
	for rows.Next() {
		var i sqlc.Command
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Status,
			&i.Source,
			&i.Inputs,
			&i.Error,
			&i.CreatedAt,
			&i.CompletedAt,
		); err != nil {
			return paging.List[model.Command]{}, fmt.Errorf("scan command: %w", err)
		}

		item, err := r.convertRowToCommand(i)
		if err != nil {
			return paging.List[model.Command]{}, fmt.Errorf("convert row to command: %w", err)
		}

		items = append(items, item)
	}

	countQuery := sq.Select("COUNT(*)").
		From("commands")

	countSQL, countArgs, err := countQuery.ToSql()
	if err != nil {
		return paging.List[model.Command]{}, fmt.Errorf("build count query: %w", err)
	}

	var count int64
	if err := db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&count); err != nil {
		return paging.List[model.Command]{}, fmt.Errorf("queries count commands: %w", err)
	}

	return paging.NewList(items, count), nil
}

func (r CommandRepository) GetCommandByStatusInProgress(ctx context.Context, sqldb db.SQLDB) (model.Command, error) {
	row, err := r.queries.CommandGetByStatusInProgress(ctx, sqldb)
	if err != nil {
		if db.IsNoRowsError(err) {
			return model.Command{}, ErrCommandNotFound
		}
		return model.Command{}, fmt.Errorf("queries get command by status in progress: %w", err)
	}

	return r.convertRowToCommand(row)
}

func (r CommandRepository) CreateCommand(ctx context.Context, db db.SQLDB, command model.Command) error {
	inputs, err := json.Marshal(command.Inputs)
	if err != nil {
		return fmt.Errorf("marshal command inputs: %w", err)
	}

	var completedAt *string
	if command.CompletedAt != nil {
		c := command.CompletedAt.Format(time.RFC3339)
		completedAt = &c
	}

	params := sqlc.CommandCreateParams{
		ID:          command.ID,
		Type:        int64(command.Type),
		Status:      int64(command.Status),
		Source:      int64(command.Source),
		Inputs:      string(inputs),
		Error:       command.Error,
		CreatedAt:   command.CreatedAt.Format(time.RFC3339),
		CompletedAt: completedAt,
	}
	if err := r.queries.CommandCreate(ctx, db, params); err != nil {
		return fmt.Errorf("queries create command: %w", err)
	}

	return nil
}

func (r CommandRepository) UpdateCommand(ctx context.Context, db db.SQLDB, params repository.UpdateCommandParams) (model.Command, error) {
	arg := sqlc.CommandUpdateParams{
		ID:             params.ID,
		SetStatus:      params.SetStatus,
		SetError:       params.SetError,
		SetCompletedAt: params.SetCompletedAt,
	}
	row, err := r.queries.CommandUpdate(ctx, db, arg)
	if err != nil {
		return model.Command{}, fmt.Errorf("queries update command: %w", err)
	}

	return r.convertRowToCommand(row)
}

func (CommandRepository) convertRowToCommand(row sqlc.Command) (model.Command, error) {
	//nolint:gosec
	inputs, err := model.UnmarshalCommandInputs(model.CommandType(row.Type), []byte(row.Inputs))
	if err != nil {
		return model.Command{}, fmt.Errorf("unmarshal command inputs: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return model.Command{}, fmt.Errorf("parse created at: %w", err)
	}

	var completedAt *time.Time
	if row.CompletedAt != nil {
		c, err := time.Parse(time.RFC3339, *row.CompletedAt)
		if err != nil {
			return model.Command{}, fmt.Errorf("parse completed at: %w", err)
		}
		completedAt = &c
	}

	//nolint:gosec
	return model.Command{
		ID:          row.ID,
		Type:        model.CommandType(row.Type),
		Status:      model.CommandStatus(row.Status),
		Source:      model.CommandSource(row.Source),
		Inputs:      inputs,
		Error:       row.Error,
		CreatedAt:   createdAt,
		CompletedAt: completedAt,
	}, nil
}

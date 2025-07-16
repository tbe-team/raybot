package alarmimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/tbe-team/raybot/internal/services/alarm"
	"github.com/tbe-team/raybot/internal/storage/db"
	"github.com/tbe-team/raybot/internal/storage/db/sqlc"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/ptr"
)

type Repository struct {
	db      db.DB
	queries *sqlc.Queries
}

func NewRepository(db db.DB, queries *sqlc.Queries) Repository {
	return Repository{
		db:      db,
		queries: queries,
	}
}

func (r Repository) ListActiveAlarms(ctx context.Context, pagingParams paging.Params) (paging.List[alarm.Alarm], error) {
	var ret paging.List[alarm.Alarm]
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		alarms, err := r.listActiveAlarms(ctx, pagingParams)
		if err != nil {
			return fmt.Errorf("failed to list active alarms: %w", err)
		}
		ret.Items = alarms
		return nil
	})

	g.Go(func() error {
		count, err := r.queries.AlarmCountActive(ctx, r.db)
		if err != nil {
			return fmt.Errorf("failed to count active alarms: %w", err)
		}
		ret.TotalItems = count
		return nil
	})

	if err := g.Wait(); err != nil {
		return paging.List[alarm.Alarm]{}, err
	}

	return ret, nil
}

//nolint:gosec
func (r Repository) listActiveAlarms(ctx context.Context, pagingParams paging.Params) ([]alarm.Alarm, error) {
	alarms, err := r.queries.AlarmListActive(ctx, r.db, sqlc.AlarmListActiveParams{
		Limit:  int64(pagingParams.Limit()),
		Offset: int64(pagingParams.Offset()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list active alarms: %w", err)
	}

	items := make([]alarm.Alarm, len(alarms))
	for i, row := range alarms {
		item, err := r.convertRowToAlarm(row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row to alarm: %w", err)
		}
		items[i] = item
	}

	return items, nil
}

func (r Repository) ListDeactiveAlarms(ctx context.Context, pagingParams paging.Params) (paging.List[alarm.Alarm], error) {
	var ret paging.List[alarm.Alarm]
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		alarms, err := r.listDeactiveAlarms(ctx, pagingParams)
		if err != nil {
			return fmt.Errorf("failed to list deactive alarms: %w", err)
		}
		ret.Items = alarms
		return nil
	})

	g.Go(func() error {
		count, err := r.queries.AlarmCountDeactive(ctx, r.db)
		if err != nil {
			return fmt.Errorf("failed to count deactive alarms: %w", err)
		}
		ret.TotalItems = count
		return nil
	})

	if err := g.Wait(); err != nil {
		return paging.List[alarm.Alarm]{}, err
	}

	return ret, nil
}

//nolint:gosec
func (r Repository) listDeactiveAlarms(ctx context.Context, pagingParams paging.Params) ([]alarm.Alarm, error) {
	alarms, err := r.queries.AlarmListDeactive(ctx, r.db, sqlc.AlarmListDeactiveParams{
		Limit:  int64(pagingParams.Limit()),
		Offset: int64(pagingParams.Offset()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list deactive alarms: %w", err)
	}

	items := make([]alarm.Alarm, len(alarms))
	for i, row := range alarms {
		item, err := r.convertRowToAlarm(row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row to alarm: %w", err)
		}
		items[i] = item
	}

	return items, nil
}

func (r Repository) UpsertActivatedAlarm(ctx context.Context, a alarm.Alarm) (alarm.Alarm, error) {
	data, err := json.Marshal(a.Data)
	if err != nil {
		return alarm.Alarm{}, fmt.Errorf("failed to marshal alarm data: %w", err)
	}

	row, err := r.queries.AlarmUpsertActivated(ctx, r.db, sqlc.AlarmUpsertActivatedParams{
		Type:        a.Type.String(),
		Data:        string(data),
		ActivatedAt: a.ActivatedAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		if db.IsUniqueViolationError(err, "alarms.type") {
			return alarm.Alarm{}, alarm.ErrAlarmAlreadyActivated
		}
		return alarm.Alarm{}, fmt.Errorf("failed to create alarm: %w", err)
	}

	return r.convertRowToAlarm(row)
}

func (r Repository) DeactivateAlarm(ctx context.Context, id int64, deactivatedAt time.Time) error {
	if err := r.queries.AlarmDeactivate(ctx, r.db, sqlc.AlarmDeactivateParams{
		ID:            id,
		DeactivatedAt: ptr.New(deactivatedAt.Format(time.RFC3339Nano)),
	}); err != nil {
		return fmt.Errorf("failed to deactivate alarm: %w", err)
	}

	return nil
}

func (r Repository) DeactivateAllAlarms(ctx context.Context) error {
	if err := r.queries.AlarmDeactivateAllActivated(ctx, r.db,
		ptr.New(time.Now().Format(time.RFC3339Nano)),
	); err != nil {
		return fmt.Errorf("failed to deactivate all activated alarms: %w", err)
	}
	return nil
}

func (r Repository) DeleteDeactivatedAlarms(ctx context.Context) error {
	if err := r.queries.AlarmDeleteDeactivated(ctx, r.db); err != nil {
		return fmt.Errorf("failed to delete deactivated alarms: %w", err)
	}

	return nil
}

func (r Repository) DeleteDeactivatedAlarmsByThreshold(ctx context.Context, threshold time.Time) error {
	if err := r.queries.AlarmDeleteDeactivatedByThreshold(ctx, r.db,
		ptr.New(threshold.Format(time.RFC3339Nano)),
	); err != nil {
		return fmt.Errorf("failed to delete deactivated alarms by threshold time: %w", err)
	}

	return nil
}

func (r Repository) convertRowToAlarm(row sqlc.Alarm) (alarm.Alarm, error) {
	ret := alarm.Alarm{
		ID:   row.ID,
		Type: alarm.AlarmType(row.Type),
	}
	var err error

	ret.Data, err = alarm.UnmarshalAlarmData(ret.Type, []byte(row.Data))
	if err != nil {
		return alarm.Alarm{}, fmt.Errorf("failed to unmarshal alarm data: %w", err)
	}

	ret.ActivatedAt, err = time.Parse(time.RFC3339Nano, row.ActivatedAt)
	if err != nil {
		return alarm.Alarm{}, fmt.Errorf("failed to parse activated at: %w", err)
	}

	if row.DeactivatedAt != nil {
		deactivatedAt, err := time.Parse(time.RFC3339Nano, *row.DeactivatedAt)
		if err != nil {
			return alarm.Alarm{}, fmt.Errorf("failed to parse deactivated at: %w", err)
		}
		ret.DeactivatedAt = &deactivatedAt
	}

	return ret, nil
}

package alarm

import (
	"context"
	"time"

	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/xerror"
)

var (
	ErrAlarmAlreadyActivated = xerror.BadRequest(nil, "alarm.alreadyActivated", "Alarm already activated")
)

type ListActiveAlarmsParams struct {
	PagingParams paging.Params `validate:"required"`
}

type ListDeactiveAlarmsParams struct {
	PagingParams paging.Params `validate:"required"`
}

type CreateAlarmParams struct {
	Data        Data      `validate:"required"`
	ActivatedAt time.Time `validate:"required"`
}

type DeactivateAlarmParams struct {
	ID            int64     `validate:"required"`
	DeactivatedAt time.Time `validate:"required"`
}

type DeleteDeactivatedAlarmsByThresholdParams struct {
	Threshold time.Time `validate:"required"`
}

type Service interface {
	ListActiveAlarms(ctx context.Context, params ListActiveAlarmsParams) (paging.List[Alarm], error)
	ListDeactiveAlarms(ctx context.Context, params ListDeactiveAlarmsParams) (paging.List[Alarm], error)
	DeleteDeactivatedAlarms(ctx context.Context) error
	DeleteDeactivatedAlarmsByThreshold(ctx context.Context, params DeleteDeactivatedAlarmsByThresholdParams) error
}

type Repository interface {
	ListActiveAlarms(ctx context.Context, pagingParams paging.Params) (paging.List[Alarm], error)
	ListDeactiveAlarms(ctx context.Context, pagingParams paging.Params) (paging.List[Alarm], error)
	UpsertActivatedAlarm(ctx context.Context, alarm Alarm) (Alarm, error)
	DeactivateAlarm(ctx context.Context, id int64, deactivatedAt time.Time) error
	DeactivateAllAlarms(ctx context.Context) error
	DeleteDeactivatedAlarms(ctx context.Context) error
	DeleteDeactivatedAlarmsByThreshold(ctx context.Context, threshold time.Time) error
}

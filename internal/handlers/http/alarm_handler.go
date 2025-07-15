package http

import (
	"context"
	"fmt"

	"github.com/tbe-team/raybot/internal/handlers/http/gen"
	"github.com/tbe-team/raybot/internal/services/alarm"
	"github.com/tbe-team/raybot/pkg/paging"
	"github.com/tbe-team/raybot/pkg/xerror"
)

type alarmHandler struct {
	alarmService alarm.Service
}

func newAlarmHandler(alarmService alarm.Service) *alarmHandler {
	return &alarmHandler{
		alarmService: alarmService,
	}
}

func (h alarmHandler) ListAlarms(ctx context.Context, req gen.ListAlarmsRequestObject) (gen.ListAlarmsResponseObject, error) {
	page := uint(1)
	pageSize := uint(10)
	if req.Params.Page != nil {
		page = *req.Params.Page
	}
	if req.Params.PageSize != nil {
		pageSize = *req.Params.PageSize
	}

	status := req.Params.Status
	pagingParams := paging.NewParams(paging.Page(page), paging.PageSize(pageSize))

	var alarms paging.List[alarm.Alarm]
	var err error

	switch status {
	case gen.ListAlarmsParamsStatusActive:
		alarms, err = h.alarmService.ListActiveAlarms(ctx, alarm.ListActiveAlarmsParams{
			PagingParams: pagingParams,
		})
		if err != nil {
			return nil, fmt.Errorf("list active alarms: %w", err)
		}

	case gen.ListAlarmsParamsStatusDeactive:
		alarms, err = h.alarmService.ListDeactiveAlarms(ctx, alarm.ListDeactiveAlarmsParams{
			PagingParams: pagingParams,
		})
		if err != nil {
			return nil, fmt.Errorf("list deactive alarms: %w", err)
		}

	default:
		return nil, xerror.ValidationFailed(nil, "invalid status")
	}

	res := make([]gen.AlarmResponse, len(alarms.Items))
	for i, alm := range alarms.Items {
		r, err := h.convertAlarmToResponse(alm)
		if err != nil {
			return nil, fmt.Errorf("convert alarm to response: %w", err)
		}
		res[i] = r
	}

	return gen.ListAlarms200JSONResponse{
		TotalItems: int(alarms.TotalItems),
		Items:      res,
	}, nil
}

func (h alarmHandler) DeleteDeactiveAlarms(ctx context.Context, _ gen.DeleteDeactiveAlarmsRequestObject) (gen.DeleteDeactiveAlarmsResponseObject, error) {
	if err := h.alarmService.DeleteDeactivatedAlarms(ctx); err != nil {
		return nil, fmt.Errorf("delete deactive alarms: %w", err)
	}

	return gen.DeleteDeactiveAlarms204Response{}, nil
}

func (h alarmHandler) convertAlarmToResponse(alm alarm.Alarm) (gen.AlarmResponse, error) {
	data, err := h.convertDataToResponse(alm.Data)
	if err != nil {
		return gen.AlarmResponse{}, fmt.Errorf("convert data to response: %w", err)
	}

	return gen.AlarmResponse{
		Id:            alm.ID,
		Type:          alm.Type.String(),
		Message:       alm.Data.Message(),
		Data:          data,
		ActivatedAt:   alm.ActivatedAt,
		DeactivatedAt: alm.DeactivatedAt,
	}, nil
}

func (alarmHandler) convertDataToResponse(data alarm.Data) (gen.AlarmData, error) {
	var res gen.AlarmData
	switch v := data.(type) {
	case alarm.DataBatteryVoltageLow:
		if err := res.FromDataBatteryVoltageLow(gen.DataBatteryVoltageLow{
			Threshold: v.Threshold,
			Voltage:   v.Voltage,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery voltage low: %w", err)
		}

	case alarm.DataBatteryVoltageHigh:
		if err := res.FromDataBatteryVoltageHigh(gen.DataBatteryVoltageHigh{
			Threshold: v.Threshold,
			Voltage:   v.Voltage,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery voltage high: %w", err)
		}

	case alarm.DataBatteryCellVoltageHigh:
		if err := res.FromDataBatteryCellVoltageHigh(gen.DataBatteryCellVoltageHigh{
			Threshold:          v.Threshold,
			CellVoltages:       v.CellVoltages,
			OverThresholdIndex: v.OverThresholdIndex,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery cell voltage high: %w", err)
		}

	case alarm.DataBatteryCellVoltageLow:
		if err := res.FromDataBatteryCellVoltageLow(gen.DataBatteryCellVoltageLow{
			Threshold:           v.Threshold,
			CellVoltages:        v.CellVoltages,
			UnderThresholdIndex: v.UnderThresholdIndex,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery cell voltage low: %w", err)
		}

	case alarm.DataBatteryCellVoltageDiff:
		if err := res.FromDataBatteryCellVoltageDiff(gen.DataBatteryCellVoltageDiff{
			Threshold:    v.Threshold,
			CellVoltages: v.CellVoltages,
			DiffIndex:    v.DiffIndex,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery cell voltage diff: %w", err)
		}

	case alarm.DataBatteryCurrentHigh:
		if err := res.FromDataBatteryCurrentHigh(gen.DataBatteryCurrentHigh{
			Threshold: v.Threshold,
			Current:   v.Current,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery current high: %w", err)
		}

	case alarm.DataBatteryTempHigh:
		if err := res.FromDataBatteryTempHigh(gen.DataBatteryTempHigh{
			Threshold: v.Threshold,
			Temp:      v.Temp,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery temp high: %w", err)
		}

	case alarm.DataBatteryPercentLow:
		if err := res.FromDataBatteryPercentLow(gen.DataBatteryPercentLow{
			Threshold: v.Threshold,
			Percent:   v.Percent,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery percent low: %w", err)
		}

	case alarm.DataBatteryHealthLow:
		if err := res.FromDataBatteryHealthLow(gen.DataBatteryHealthLow{
			Threshold: v.Threshold,
			Health:    v.Health,
		}); err != nil {
			return gen.AlarmData{}, fmt.Errorf("from data battery health low: %w", err)
		}

	default:
		return gen.AlarmData{}, fmt.Errorf("unknown data type: %T", v)
	}

	return res, nil
}

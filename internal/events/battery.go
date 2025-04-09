package events

import "github.com/maniartech/signals"

type UpdateBatteryChargeSettingEvent struct {
	CurrentLimit uint16
	Enable       bool
}

type UpdateBatteryDischargeSettingEvent struct {
	CurrentLimit uint16
	Enable       bool
}

var UpdateBatteryChargeSettingSignal = signals.New[UpdateBatteryChargeSettingEvent]()
var UpdateBatteryDischargeSettingSignal = signals.New[UpdateBatteryDischargeSettingEvent]()

export const AlarmTypeValues = [
  'battery_voltage_low',
  'battery_voltage_high',
  'battery_cell_voltage_high',
  'battery_cell_voltage_low',
  'battery_cell_voltage_diff',
  'battery_current_high',
  'battery_temp_high',
  'battery_percent_low',
  'battery_health_low',
] as const

export type AlarmType = typeof AlarmTypeValues[number]

export interface DataBatteryVoltageLow {
  threshold: number
  voltage: number
}

export interface DataBatteryVoltageHigh {
  threshold: number
  voltage: number
}

export interface DataBatteryCellVoltageHigh {
  threshold: number
  cellVoltages: number[]
  overThresholdIndex: number[]
}

export interface DataBatteryCellVoltageLow {
  threshold: number
  cellVoltages: number[]
  overThresholdIndex: number[]
}

export interface DataBatteryCellVoltageDiff {
  threshold: number
  cellVoltages: number[]
  diffIndex: number[]
}

export interface DataBatteryCurrentHigh {
  threshold: number
  current: number
}

export interface DataBatteryTempHigh {
  threshold: number
  temp: number
}

export interface DataBatteryPercentLow {
  threshold: number
  percent: number
}

export interface DataBatteryHealthLow {
  threshold: number
  health: number
}

export interface AlarmData {
  battery_voltage_low: DataBatteryVoltageLow
  battery_voltage_high: DataBatteryVoltageHigh
  battery_cell_voltage_high: DataBatteryCellVoltageHigh
  battery_cell_voltage_low: DataBatteryCellVoltageLow
  battery_cell_voltage_diff: DataBatteryCellVoltageDiff
  battery_current_high: DataBatteryCurrentHigh
  battery_temp_high: DataBatteryTempHigh
  battery_percent_low: DataBatteryPercentLow
  battery_health_low: DataBatteryHealthLow
}

export interface Alarm<T extends AlarmType = AlarmType> {
  id: number
  type: T
  message: string
  data: AlarmData[T]
  activatedAt: string
  deactivatedAt?: string
}

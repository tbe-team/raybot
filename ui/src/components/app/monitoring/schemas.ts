import z from 'zod'

export const monitoringConfigSchema = z.object({
  voltageLow: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0),
  }),
  voltageHigh: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0),
  }),
  cellVoltageHigh: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0),
  }),
  cellVoltageLow: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0),
  }),
  cellVoltageDiff: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0),
  }),
  currentHigh: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0),
  }),
  tempHigh: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0),
  }),
  percentLow: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0).max(100),
  }),
  healthLow: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0).max(100),
  }),
})

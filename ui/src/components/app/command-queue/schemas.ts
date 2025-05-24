import { z } from 'zod'

export const createCommandSchema = z.discriminatedUnion('type', [
  z.object({
    type: z.literal('STOP_MOVEMENT'),
    inputs: z.object({}).default({}),
  }),
  z.object({
    type: z.literal('MOVE_TO'),
    inputs: z.object({
      location: z.string(),
      direction: z.union([z.literal('FORWARD'), z.literal('BACKWARD')]),
      motorSpeed: z.number().min(0).max(100),
    }),
  }),
  z.object({
    type: z.literal('MOVE_FORWARD'),
    inputs: z.object({
      motorSpeed: z.number().min(0).max(100),
    }),
  }),
  z.object({
    type: z.literal('MOVE_BACKWARD'),
    inputs: z.object({
      motorSpeed: z.number().min(0).max(100),
    }),
  }),
  z.object({
    type: z.literal('CARGO_OPEN'),
    inputs: z.object({
      motorSpeed: z.number().min(0).max(100),
    }),
  }),
  z.object({
    type: z.literal('CARGO_CLOSE'),
    inputs: z.object({
      motorSpeed: z.number().min(0).max(100),
    }),
  }),
  z.object({
    type: z.literal('CARGO_LIFT'),
    inputs: z.object({
      motorSpeed: z.number().min(0).max(100),
      position: z.number().min(0),
    }),
  }),
  z.object({
    type: z.literal('CARGO_LOWER'),
    inputs: z.object({
      motorSpeed: z.number().min(0).max(100),
      bottomObstacleTracking: z.object({
        enterDistance: z.number().min(0),
        exitDistance: z.number().min(0),
      }),
      position: z.number().min(0),
    }),
  }),
  z.object({
    type: z.literal('CARGO_CHECK_QR'),
    inputs: z.object({
      qrCode: z.string(),
    }),
  }),
  z.object({
    type: z.literal('SCAN_LOCATION'),
    inputs: z.object({}).default({}),
  }),
  z.object({
    type: z.literal('WAIT'),
    inputs: z.object({
      durationMs: z.number(),
    }),
  }),
])

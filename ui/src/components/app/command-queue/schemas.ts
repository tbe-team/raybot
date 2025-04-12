import { z } from 'zod'

export const createCommandSchema = z.discriminatedUnion('type', [
  z.object({
    type: z.literal('STOP'),
    inputs: z.object({}),
  }),
  z.object({
    type: z.literal('MOVE_TO'),
    inputs: z.object({
      location: z.string(),
    }),
  }),
  z.object({
    type: z.literal('MOVE_FORWARD'),
    inputs: z.object({}),
  }),
  z.object({
    type: z.literal('MOVE_BACKWARD'),
    inputs: z.object({}),
  }),
  z.object({
    type: z.literal('CARGO_OPEN'),
    inputs: z.object({}),
  }),
  z.object({
    type: z.literal('CARGO_CLOSE'),
    inputs: z.object({}),
  }),
  z.object({
    type: z.literal('CARGO_LIFT'),
    inputs: z.object({}),
  }),
  z.object({
    type: z.literal('CARGO_LOWER'),
    inputs: z.object({}),
  }),
  z.object({
    type: z.literal('CARGO_CHECK_QR'),
    inputs: z.object({
      qrCode: z.string(),
    }),
  }),
])

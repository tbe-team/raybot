import * as z from 'zod'

export const systemConfigSchema = z.object({
  grpc: z.object({
    server: z.object({
      enable: z.boolean(),
    }),
    cloud: z.object({
      address: z.string(),
    }),
  }),
  http: z.object({
    enableSwagger: z.boolean(),
  }),
  log: z.object({
    level: z.union([z.literal('debug'), z.literal('info'), z.literal('warn'), z.literal('error')]),
    format: z.union([z.literal('json'), z.literal('text')]),
    addSource: z.boolean(),
  }),
  pic: z.object({
    serial: z.object({
      port: z.string(),
      baudRate: z.number().int(),
      dataBits: z.union([z.literal(5), z.literal(6), z.literal(7), z.literal(8)]),
      stopBits: z.union([z.literal(1), z.literal(1.5), z.literal(2)]),
      parity: z.union([z.literal('none'), z.literal('even'), z.literal('odd')]),
      readTimeout: z.number().min(0),
    }),
  }),
})

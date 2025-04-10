import type { Command } from '@/types/command'

const fakeCommands: Command[] = [
  {
    id: 1,
    type: 'STOP',
    status: 'SUCCEEDED',
    source: 'MANUAL',
    inputs: {},
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: 2,
    type: 'MOVE_FORWARD',
    status: 'PROCESSING',
    source: 'APP',
    inputs: {},
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: 3,
    type: 'MOVE_BACKWARD',
    status: 'FAILED',
    source: 'MANUAL',
    inputs: {},
    error: 'Failed to move backward',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: 4,
    type: 'MOVE_TO',
    status: 'SUCCEEDED',
    source: 'APP',
    inputs: { location: 'Location A' },
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: 5,
    type: 'CARGO_OPEN',
    status: 'PROCESSING',
    source: 'MANUAL',
    inputs: {},
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
]

export default fakeCommands

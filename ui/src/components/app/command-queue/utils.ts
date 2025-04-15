import type { CommandType } from '@/types/command'
import type { LucideIcon } from 'lucide-vue-next'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import {
  ArrowDown,
  ArrowUp,
  MapPin,
  Package,
  QrCode,
  StopCircle,
} from 'lucide-vue-next'

dayjs.extend(relativeTime)

const commandIcons: Record<CommandType, LucideIcon> = {
  STOP_MOVEMENT: StopCircle,
  MOVE_FORWARD: ArrowUp,
  MOVE_BACKWARD: ArrowDown,
  MOVE_TO: MapPin,
  CARGO_OPEN: Package,
  CARGO_CLOSE: Package,
  CARGO_LIFT: Package,
  CARGO_LOWER: Package,
  CARGO_CHECK_QR: QrCode,
}

export function getCommandIcon(type: CommandType) {
  return commandIcons[type]
}

export function getCommandName(type: CommandType) {
  return type
    .replace(/_/g, ' ')
    .toLowerCase()
    .replace(/\b\w/g, l => l.toUpperCase())
}

export function timeSince(dateStr: string) {
  const date = dayjs(dateStr)
  return date.isValid() ? date.fromNow() : 'N/A'
}

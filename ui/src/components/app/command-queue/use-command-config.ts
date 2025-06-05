import type {
  CargoCloseInputs,
  CargoLiftInputs,
  CargoLowerInputs,
  CargoOpenInputs,
  CommandInputMap,
  CommandType,
  MoveBackwardInputs,
  MoveForwardInputs,
  MoveToInputs,
} from '@/types/command'
import { useLocalStorage } from '@vueuse/core'

export interface CommandConfig {
  moveTo: Omit<MoveToInputs, 'location'>
  moveForward: MoveForwardInputs
  moveBackward: MoveBackwardInputs
  cargoOpen: CargoOpenInputs
  cargoClose: CargoCloseInputs
  cargoLift: CargoLiftInputs
  cargoLower: CargoLowerInputs
}

const STORAGE_KEY = 'command-config'
const DEFAULT_COMMAND_CONFIG: CommandConfig = {
  moveTo: {
    direction: 'FORWARD',
    motorSpeed: 80,
  },
  moveForward: {
    motorSpeed: 80,
  },
  moveBackward: {
    motorSpeed: 80,
  },
  cargoOpen: {
    motorSpeed: 80,
  },
  cargoClose: {
    motorSpeed: 80,
  },
  cargoLift: {
    motorSpeed: 80,
    position: 20,
  },
  cargoLower: {
    motorSpeed: 80,
    position: 100,
  },
}

export function useCommandConfig() {
  const state = useLocalStorage<CommandConfig>(STORAGE_KEY, DEFAULT_COMMAND_CONFIG)

  const commandConfig = computed(() => state.value)

  const updateHandlers: {
    [K in CommandType]: (inputs: CommandInputMap[K]) => void
  } = {
    STOP_MOVEMENT: () => {},
    MOVE_BACKWARD: (inputs) => {
      state.value = {
        ...state.value,
        moveBackward: {
          motorSpeed: inputs.motorSpeed,
        },
      }
    },
    MOVE_FORWARD: (inputs) => {
      state.value = {
        ...state.value,
        moveForward: {
          motorSpeed: inputs.motorSpeed,
        },
      }
    },
    MOVE_TO: (inputs) => {
      state.value = {
        ...state.value,
        moveTo: {
          direction: inputs.direction,
          motorSpeed: inputs.motorSpeed,
        },
      }
    },
    CARGO_OPEN: (inputs) => {
      state.value = {
        ...state.value,
        cargoOpen: {
          motorSpeed: inputs.motorSpeed,
        },
      }
    },
    CARGO_CLOSE: (inputs) => {
      state.value = {
        ...state.value,
        cargoClose: {
          motorSpeed: inputs.motorSpeed,
        },
      }
    },
    CARGO_LIFT: (inputs) => {
      state.value = {
        ...state.value,
        cargoLift: {
          motorSpeed: inputs.motorSpeed,
          position: inputs.position,
        },
      }
    },
    CARGO_LOWER: (inputs) => {
      state.value = {
        ...state.value,
        cargoLower: {
          motorSpeed: inputs.motorSpeed,
          position: inputs.position,
        },
      }
    },
    CARGO_CHECK_QR: () => {},
    SCAN_LOCATION: () => {},
    WAIT: () => {},
  }

  function updateCommandConfigFromInputs(type: CommandType, inputs: CommandInputMap[CommandType]) {
    const handler = updateHandlers[type]
    if (handler) {
      handler(inputs as never)
    }
  }

  return {
    commandConfig,
    updateCommandConfigFromInputs,
  }
}

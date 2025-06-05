export interface CargoLiftConfig {
  stableReadCount: number
}

export interface ObstacleTracking {
  enterDistance: number
  exitDistance: number
}

export interface CargoLowerConfig {
  stableReadCount: number
  bottomObstacleTracking: ObstacleTracking
}

export interface CommandConfig {
  cargoLift: CargoLiftConfig
  cargoLower: CargoLowerConfig
}

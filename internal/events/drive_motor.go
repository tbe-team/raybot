package events

import "github.com/tbe-team/raybot/internal/services/drivemotor"

const (
	DriveMotorUpdatedTopic = "drive_motor_updated"
)

type DriveMotorStateUpdatedEvent struct {
	Direction drivemotor.Direction `json:"direction"`
	Speed     uint8                `json:"speed"`
	IsRunning bool                 `json:"is_running"`
	Enabled   bool                 `json:"enabled"`
}

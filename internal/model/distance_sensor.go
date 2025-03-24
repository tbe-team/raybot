package model

import "time"

type DistanceSensor struct {
	FrontDistance uint16
	BackDistance  uint16
	DownDistance  uint16
	UpdatedAt     time.Time
}

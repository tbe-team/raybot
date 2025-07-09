package system

import "time"

type Info struct {
	LocalIP string
	// CPU usage in percentage
	CPUUsage float64
	// Memory usage in percentage
	MemoryUsage float64
	// Total memory in MB
	TotalMemory uint64
	Uptime      time.Duration
}

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	StatusNormal Status = "NORMAL"
	StatusError  Status = "ERROR"
)

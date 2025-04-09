package peripheral

import "context"

type Service interface {
	ListAvailableSerialPorts(ctx context.Context) ([]SerialPort, error)
}

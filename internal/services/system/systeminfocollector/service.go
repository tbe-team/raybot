package systeminfocollector

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/tbe-team/raybot/internal/services/system"
)

const (
	sysMetricsCollectInterval  = 10 * time.Second
	networkInfoCollectInterval = 1 * time.Minute
)

type service struct {
	log            *slog.Logger
	systemInfoRepo system.Repository
	stopCh         chan struct{}
}

func NewService(log *slog.Logger, systemInfoRepo system.Repository) system.SysInfoCollectorService {
	return &service{
		log:            log,
		systemInfoRepo: systemInfoRepo,
		stopCh:         make(chan struct{}),
	}
}

func (s service) Run(ctx context.Context) {
	go func() {
		if err := s.collectSystemMetrics(ctx); err != nil {
			s.log.Error("failed to collect system metrics", "error", err)
		}

		for {
			select {
			case <-ctx.Done():
				return

			case <-s.stopCh:
				return

			case <-time.After(sysMetricsCollectInterval):
				if err := s.collectSystemMetrics(ctx); err != nil {
					s.log.Error("failed to collect system metrics", "error", err)
				}
			}
		}
	}()

	go func() {
		if err := s.collectNetworkInfo(ctx); err != nil {
			s.log.Error("failed to collect network info", "error", err)
		}

		for {
			select {
			case <-ctx.Done():
				return

			case <-s.stopCh:
				return

			case <-time.After(networkInfoCollectInterval):
				if err := s.collectNetworkInfo(ctx); err != nil {
					s.log.Error("failed to collect network info", "error", err)
				}
			}
		}
	}()
}

func (s service) Stop() {
	close(s.stopCh)
}

func (s service) collectSystemMetrics(ctx context.Context) error {
	params := system.UpdateInfoParams{}

	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return fmt.Errorf("failed to get CPU usage: %w", err)
	}

	if len(cpuUsage) > 0 {
		params.CPUUsage = roundTo2Decimal(cpuUsage[0])
		params.SetCPUUsage = true
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to get memory usage: %w", err)
	}

	params.MemoryUsage = roundTo2Decimal(memory.UsedPercent)
	params.SetMemoryUsage = true
	// Convert to MB
	params.TotalMemory = memory.Total / 1024 / 1024
	params.SetTotalMemory = true

	uptime, err := host.Uptime()
	if err != nil {
		return fmt.Errorf("failed to get uptime: %w", err)
	}

	//nolint:gosec
	params.Uptime = time.Duration(uptime) * time.Second
	params.SetUptime = true

	if err := s.systemInfoRepo.UpdateInfo(ctx, params); err != nil {
		return fmt.Errorf("failed to update system info: %w", err)
	}

	return nil
}

func (s service) collectNetworkInfo(ctx context.Context) error {
	networkInfo := system.UpdateInfoParams{}
	networkInfo.LocalIP = getLocalIP()
	networkInfo.SetLocalIP = true

	if err := s.systemInfoRepo.UpdateInfo(ctx, networkInfo); err != nil {
		return fmt.Errorf("failed to update network info: %w", err)
	}

	return nil
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func roundTo2Decimal(f float64) float64 {
	return math.Round(f*100) / 100
}

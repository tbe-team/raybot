package monitoringimpl

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/events"
	"github.com/tbe-team/raybot/internal/services/alarm"
	"github.com/tbe-team/raybot/internal/services/battery"
	configsvc "github.com/tbe-team/raybot/internal/services/config"
	"github.com/tbe-team/raybot/pkg/eventbus"
)

type Service struct {
	log             *slog.Logger
	eventSubscriber eventbus.Subscriber
	alarmRepo       alarm.Repository
	batteryRepo     battery.BatteryStateRepository
	configService   configsvc.Service

	stopCh chan struct{}
	wg     sync.WaitGroup
}

func NewService(
	log *slog.Logger,
	eventSubscriber eventbus.Subscriber,
	alarmRepo alarm.Repository,
	batteryRepo battery.BatteryStateRepository,
	configService configsvc.Service,
) *Service {
	return &Service{
		log:             log,
		eventSubscriber: eventSubscriber,
		alarmRepo:       alarmRepo,
		batteryRepo:     batteryRepo,
		configService:   configService,
		stopCh:          make(chan struct{}),
	}
}

func (s *Service) Start(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			err := s.startMonitorBattery(ctx)
			if err == nil {
				return
			}

			s.log.Error("failed to start monitor battery, retrying in 5 seconds", slog.Any("error", err))
			select {
			case <-ctx.Done():
				return

			case <-s.stopCh:
				return

			case <-time.After(5 * time.Second):
				continue
			}
		}
	}()
}

func (s *Service) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

func (s *Service) startMonitorBattery(ctx context.Context) error {
	cfg, err := s.configService.GetBatteryMonitoringConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get battery monitoring config: %w", err)
	}

	// check battery alarms on init
	batteryState, err := s.batteryRepo.GetBatteryState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get battery state: %w", err)
	}
	s.checkBatteryAlarms(ctx, batteryState, cfg)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel to unsubscribe from the event

	s.eventSubscriber.Subscribe(ctx, events.BatteryUpdatedTopic, func(msg *eventbus.Message) {
		ev, ok := msg.Payload.(events.BatteryUpdatedEvent)
		if !ok {
			s.log.Error("invalid battery updated event", slog.Any("event", msg.Payload))
			return
		}

		s.checkBatteryAlarms(ctx, ev.BatteryState, cfg)
	})

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-s.stopCh:
			return nil
		}
	}
}

func (s *Service) checkBatteryAlarms(ctx context.Context, state battery.BatteryState, cfg config.BatteryMonitoring) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.checkVoltageAlarms(ctx, state, cfg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.checkCellVoltageAlarms(ctx, state, cfg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.checkCurrentAlarms(ctx, state, cfg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.checkTemperatureAlarms(ctx, state, cfg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.checkPercentAlarms(ctx, state, cfg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.checkHealthAlarms(ctx, state, cfg)
	}()

	wg.Wait()
}

func (s *Service) checkVoltageAlarms(ctx context.Context, state battery.BatteryState, config config.BatteryMonitoring) {
	voltage := float64(state.Voltage)

	if config.BatteryVoltageLow.Enable && voltage < config.BatteryVoltageLow.Threshold {
		s.createAlarm(ctx, alarm.DataBatteryVoltageLow{
			Threshold: config.BatteryVoltageLow.Threshold,
			Voltage:   voltage,
		})
	}

	if config.BatteryVoltageHigh.Enable && voltage > config.BatteryVoltageHigh.Threshold {
		s.createAlarm(ctx, alarm.DataBatteryVoltageHigh{
			Threshold: config.BatteryVoltageHigh.Threshold,
			Voltage:   voltage,
		})
	}
}

func (s *Service) checkCellVoltageAlarms(ctx context.Context, state battery.BatteryState, config config.BatteryMonitoring) {
	if len(state.CellVoltages) == 0 {
		return
	}

	cellVoltages := make([]float64, len(state.CellVoltages))
	for i, v := range state.CellVoltages {
		cellVoltages[i] = float64(v)
	}

	if config.BatteryCellVoltageHigh.Enable {
		overThresholdIndex := []int{}
		for i, voltage := range cellVoltages {
			if voltage > config.BatteryCellVoltageHigh.Threshold {
				overThresholdIndex = append(overThresholdIndex, i)
			}
		}
		if len(overThresholdIndex) > 0 {
			s.createAlarm(ctx, alarm.DataBatteryCellVoltageHigh{
				Threshold:          config.BatteryCellVoltageHigh.Threshold,
				CellVoltages:       cellVoltages,
				OverThresholdIndex: overThresholdIndex,
			})
		}
	}

	if config.BatteryCellVoltageLow.Enable {
		underThresholdIndex := []int{}
		for i, voltage := range cellVoltages {
			if voltage < config.BatteryCellVoltageLow.Threshold {
				underThresholdIndex = append(underThresholdIndex, i)
			}
		}
		if len(underThresholdIndex) > 0 {
			s.createAlarm(ctx, alarm.DataBatteryCellVoltageLow{
				Threshold:           config.BatteryCellVoltageLow.Threshold,
				CellVoltages:        cellVoltages,
				UnderThresholdIndex: underThresholdIndex,
			})
		}
	}

	if config.BatteryCellVoltageDiff.Enable {
		minVoltage, maxVoltage := cellVoltages[0], cellVoltages[0]
		for _, voltage := range cellVoltages {
			minVoltage = math.Min(minVoltage, voltage)
			maxVoltage = math.Max(maxVoltage, voltage)
		}

		voltageDiff := maxVoltage - minVoltage
		if voltageDiff > config.BatteryCellVoltageDiff.Threshold {
			diffIndex := []int{}
			for i := range cellVoltages {
				diffIndex = append(diffIndex, i)
			}
			s.createAlarm(ctx, alarm.DataBatteryCellVoltageDiff{
				Threshold:    config.BatteryCellVoltageDiff.Threshold,
				CellVoltages: cellVoltages,
				DiffIndex:    diffIndex,
			})
		}
	}
}

func (s *Service) checkCurrentAlarms(ctx context.Context, state battery.BatteryState, config config.BatteryMonitoring) {
	current := float64(state.Current)

	if config.BatteryCurrentHigh.Enable && current > config.BatteryCurrentHigh.Threshold {
		s.createAlarm(ctx, alarm.DataBatteryCurrentHigh{
			Threshold: config.BatteryCurrentHigh.Threshold,
			Current:   current,
		})
	}
}

func (s *Service) checkTemperatureAlarms(ctx context.Context, state battery.BatteryState, config config.BatteryMonitoring) {
	temp := float64(state.Temp)

	if config.BatteryTempHigh.Enable && temp > config.BatteryTempHigh.Threshold {
		s.createAlarm(ctx, alarm.DataBatteryTempHigh{
			Threshold: config.BatteryTempHigh.Threshold,
			Temp:      temp,
		})
	}
}

func (s *Service) checkPercentAlarms(ctx context.Context, state battery.BatteryState, config config.BatteryMonitoring) {
	percent := float64(state.Percent)

	if config.BatteryPercentLow.Enable && percent < config.BatteryPercentLow.Threshold {
		s.createAlarm(ctx, alarm.DataBatteryPercentLow{
			Threshold: config.BatteryPercentLow.Threshold,
			Percent:   percent,
		})
	}
}

func (s *Service) checkHealthAlarms(ctx context.Context, state battery.BatteryState, config config.BatteryMonitoring) {
	health := float64(state.Health)

	if config.BatteryHealthLow.Enable && health < config.BatteryHealthLow.Threshold {
		s.createAlarm(ctx, alarm.DataBatteryHealthLow{
			Threshold: config.BatteryHealthLow.Threshold,
			Health:    health,
		})
	}
}

func (s *Service) createAlarm(ctx context.Context, data alarm.Data) {
	_, err := s.alarmRepo.UpsertActivatedAlarm(ctx, alarm.Alarm{
		Type:        data.AlarmType(),
		Data:        data,
		ActivatedAt: time.Now(),
	})
	if err != nil {
		s.log.Error("failed to create battery alarm", slog.Any("error", err))
		return
	}
}

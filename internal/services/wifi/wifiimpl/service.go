package wifiimpl

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/tbe-team/raybot/internal/config"
	"github.com/tbe-team/raybot/internal/services/wifi"
)

const (
	apName  = "raybot-ap-mode"
	staName = "raybot-sta-mode"
)

type service struct {
	cfg config.Wifi
	log *slog.Logger
}

func NewService(cfg config.Wifi, log *slog.Logger) wifi.Service {
	return &service{
		cfg: cfg,
		log: log.With("service", "wifi"),
	}
}

func (s service) Run(_ context.Context) error {
	if err := s.initWifi(); err != nil {
		return fmt.Errorf("init wifi: %w", err)
	}

	return nil
}

func (s service) initWifi() error {
	// allow bypass if both AP and STA are disabled
	if !s.cfg.AP.Enable && !s.cfg.STA.Enable {
		return nil
	}

	currentMode, err := getCurrentMode()
	if err != nil {
		return fmt.Errorf("get current mode: %w", err)
	}

	if currentMode == wifiModeAP && s.cfg.STA.Enable {
		if err := downConnection(apName); err != nil {
			return fmt.Errorf("down connection: %w", err)
		}
	}

	if currentMode == wifiModeSTA && s.cfg.AP.Enable {
		if err := downConnection(staName); err != nil {
			return fmt.Errorf("down connection: %w", err)
		}
	}

	if s.cfg.AP.Enable {
		return s.initAPMode()
	}

	if s.cfg.STA.Enable {
		if err := s.initSTAMode(); err != nil {
			s.log.Error("failed to initialize STA mode, falling back to AP mode", slog.Any("error", err))
			return s.initAPMode()
		}
	}

	return nil
}

func (s service) initAPMode() error {
	s.log.Info("initializing AP mode", slog.Any("config", s.cfg.AP))

	if err := deleteConnection(apName); err != nil {
		return fmt.Errorf("delete connection: %w", err)
	}

	//nolint:gosec
	cmd := exec.Command(
		"nmcli", "con", "add",
		"type", "wifi",
		"con-name", apName,
		"autoconnect", "no",
		"ssid", s.cfg.AP.SSID,
		"802-11-wireless.mode", "ap",
		"802-11-wireless.band", "bg",
		"802-11-wireless.channel", "7",
		"ipv4.addresses", s.cfg.AP.IP+"/24",
		"ipv4.method", "shared",
		"wifi-sec.key-mgmt", "wpa-psk",
		"wifi-sec.psk", s.cfg.AP.Password,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("create connection error: %v\nOutput: %s", err, out)
	}

	if err := upConnection(apName); err != nil {
		return fmt.Errorf("up connection: %w", err)
	}

	s.log.Info("AP mode initialized successfully", slog.Any("config", s.cfg.AP))

	return nil
}

func (s service) initSTAMode() error {
	s.log.Info("initializing STA mode", slog.Any("config", s.cfg.STA))

	if err := deleteConnection(staName); err != nil {
		return fmt.Errorf("delete connection: %w", err)
	}

	//nolint:gosec
	cmd := exec.Command(
		"nmcli", "con", "add",
		"type", "wifi",
		"con-name", staName,
		"autoconnect", "no",
		"ssid", s.cfg.STA.SSID,
		"wifi-sec.key-mgmt", "wpa-psk",
		"wifi-sec.psk", s.cfg.STA.Password,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("connect to wifi error: %v\nOutput: %s", err, out)
	}

	if err := upConnection(staName); err != nil {
		return fmt.Errorf("up connection: %w", err)
	}

	s.log.Info("STA mode initialized successfully", slog.Any("config", s.cfg.STA))

	return nil
}

func upConnection(conName string) error {
	cmd := exec.Command("nmcli", "con", "up", conName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("up connection error: %v\nOutput: %s", err, out)
	}

	return nil
}

func downConnection(conName string) error {
	cmd := exec.Command("nmcli", "con", "down", conName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("down connection error: %v\nOutput: %s", err, out)
	}

	return nil
}

func deleteConnection(conName string) error {
	cmd := exec.Command("nmcli", "con", "delete", conName)
	out, err := cmd.CombinedOutput()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 10 { // nmcli code 10 means not found
			return fmt.Errorf("delete connection error (exit code %d): %v\nOutput: %s", exitErr.ExitCode(), err, out)
		}
	}

	return nil
}

func getCurrentMode() (wifiMode, error) {
	cmd := exec.Command(
		"nmcli", "-t", "-f", "NAME", "con", "show", "--active",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return wifiModeUnknown, fmt.Errorf("get current mode error: %v\nOutput: %s", err, out)
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, apName) {
			return wifiModeAP, nil
		}
		if strings.Contains(line, staName) {
			return wifiModeSTA, nil
		}
	}

	return wifiModeUnknown, nil
}

type wifiMode uint8

const (
	wifiModeUnknown wifiMode = iota
	wifiModeAP
	wifiModeSTA
)

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

const apName = "raybot-ap-mode"

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

	// check if AP mode is already up
	check := exec.Command("nmcli", "-t", "-f", "NAME,DEVICE", "con", "show", "--active")
	out, _ := check.CombinedOutput()
	if strings.Contains(string(out), apName) {
		s.log.Info("AP mode already active, skipping setup")
		return nil
	}

	// delete existing connection
	cmd := exec.Command("nmcli", "con", "delete", apName)
	out, err := cmd.CombinedOutput()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 10 { // nmcli code 10 means not found
			return fmt.Errorf("delete connection error (exit code %d): %v\nOutput: %s", exitErr.ExitCode(), err, out)
		}
	}

	//nolint:gosec
	cmd = exec.Command(
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
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("create connection error: %v\nOutput: %s", err, out)
	}

	cmd = exec.Command("nmcli", "con", "up", apName)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("up connection error: %v\nOutput: %s", err, out)
	}

	s.log.Info("AP mode initialized successfully", slog.Any("config", s.cfg.AP))

	return nil
}

func (s service) initSTAMode() error {
	// check current connection first
	checkCmd := exec.Command("nmcli", "-t", "-f", "ACTIVE,SSID", "dev", "wifi")
	out, _ := checkCmd.CombinedOutput()
	lines := string(out)
	if strings.Contains(lines, "yes:"+s.cfg.STA.SSID) {
		s.log.Info("already connected to wifi network, skipping connect")
		return nil
	}

	//nolint:gosec
	cmd := exec.Command("nmcli", "device", "wifi",
		"connect", s.cfg.STA.SSID,
		"password", s.cfg.STA.Password,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("connect to wifi error: %v\nOutput: %s", err, out)
	}

	return nil
}

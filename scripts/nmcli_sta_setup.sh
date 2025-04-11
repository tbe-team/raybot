#!/bin/bash

set -e

if [[ $EUID -ne 0 ]]; then
  echo "⚠️ Must be run as root"
  exit 1
fi

if [[ $# -lt 2 ]]; then
  echo "Usage: sudo $0 <SSID> <PASSWORD>"
  exit 1
fi

SSID="$1"
PASS="$2"
CON_NAME="sta-mode"

echo "📶 Switching to WiFi Client mode with SSID: $SSID"

# Delete if exists
nmcli con delete "$CON_NAME" 2>/dev/null || true

# Connect and save profile
nmcli device wifi connect "$SSID" password "$PASS" ifname wlan0 name "$CON_NAME"

echo "✅ Client mode is now active and connected to $SSID"

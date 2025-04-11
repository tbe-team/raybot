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
CON_NAME="ap-mode"

echo "📡 Creating AP mode profile with SSID: $SSID"

# Delete if exists
nmcli con delete "$CON_NAME" 2>/dev/null || true

# Create AP connection
nmcli connection add type wifi ifname wlan0 con-name "$CON_NAME" autoconnect no ssid "$SSID"
nmcli connection modify "$CON_NAME" 802-11-wireless.mode ap
nmcli connection modify "$CON_NAME" 802-11-wireless.band bg
nmcli connection modify "$CON_NAME" 802-11-wireless.channel 7
nmcli connection modify "$CON_NAME" ipv4.method shared
nmcli connection modify "$CON_NAME" wifi-sec.key-mgmt wpa-psk
nmcli connection modify "$CON_NAME" wifi-sec.psk "$PASS"

# Bring it up
nmcli con up "$CON_NAME"

echo "✅ AP mode is now active on SSID: $SSID"

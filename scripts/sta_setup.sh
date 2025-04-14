#!/bin/bash

set -e

if [ "$EUID" -ne 0 ]; then
  echo "⚠️  Must be run as root"
  exit 1
fi

echo "↩️ Switching back to WiFi Client mode"

# Stop AP services
systemctl stop hostapd || true
systemctl stop dnsmasq || true
systemctl disable hostapd || true
systemctl disable dnsmasq || true

# Remove AP config
rm -f /etc/network/interfaces.d/wlan0-ap
rm -f /etc/hostapd/hostapd.conf
rm -f /etc/dnsmasq.conf

# Remove hostapd line from /etc/default/hostapd
sed -i 's|^DAEMON_CONF=.*|#DAEMON_CONF=""|' /etc/default/hostapd

# Unblock dhcpcd managing wlan0
sed -i '/denyinterfaces wlan0/d' /etc/dhcpcd.conf

systemctl restart dhcpcd

echo "✅ WiFi client mode restored. Now edit /etc/wpa_supplicant/wpa_supplicant.conf to set your network."
echo "🔁 Reboot for cleanest result."

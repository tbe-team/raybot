#!/bin/bash
# ref: https://gist.github.com/Lewiscowles1986/fecd4de0b45b2029c390

set -e

# === CONFIG AREA === #
APSSID="$1"
APPASS="$2"

STATIC_IP="192.168.1.11"
DHCP_START="192.168.1.20"
DHCP_END="192.168.1.22"
DHCP_LEASE="12h"
CHANNEL="7"
# ==================== #

if [ "$EUID" -ne 0 ]; then
  echo "⚠️  Must be run as root"
  exit 1
fi

if [[ $# -lt 2 ]]; then
  echo "Usage: sudo $0 <AP_SSID> <AP_PASSWORD>"
  exit 1
fi

echo "🚀 Switching to Access Point mode with SSID: $APSSID"

apt-get remove --purge hostapd -yqq
apt-get update -yqq
apt-get install hostapd dnsmasq -yqq

# Dnsmasq config
cat > /etc/dnsmasq.conf <<EOF
interface=wlan0
dhcp-range=$DHCP_START,$DHCP_END,$NETMASK,$DHCP_LEASE
EOF

# Hostapd config
cat > /etc/hostapd/hostapd.conf <<EOF
interface=wlan0
hw_mode=g
channel=$CHANNEL
auth_algs=1
wpa=2
wpa_key_mgmt=WPA-PSK
wpa_passphrase=$APPASS
ssid=$APSSID
rsn_pairwise=CCMP
ieee80211n=1
wmm_enabled=1
ht_capab=[HT40][SHORT-GI-20][DSSS_CCK-40]
EOF

# Interface static IP
cat > /etc/network/interfaces.d/wlan0-ap <<EOF
auto wlan0
iface wlan0 inet static
    address $STATIC_IP
    netmask 255.255.255.0
    network 192.168.1.0
    broadcast 192.168.1.255
EOF

# Block dhcpcd from managing wlan0
grep -q "denyinterfaces wlan0" /etc/dhcpcd.conf || echo "denyinterfaces wlan0" >> /etc/dhcpcd.conf

# Set hostapd daemon config
sed -i 's|^#DAEMON_CONF="".*|DAEMON_CONF="/etc/hostapd/hostapd.conf"|' /etc/default/hostapd

# Enable services
systemctl unmask hostapd
systemctl enable hostapd
systemctl enable dnsmasq

systemctl restart dhcpcd
systemctl start hostapd
systemctl start dnsmasq

echo "✅ AP mode enabled! SSID: $APSSID | Password: $APPASS"
echo "🔁 Please reboot to fully apply the network config"

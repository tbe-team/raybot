package main

import (
	"fmt"
	"log"
	"net"
)

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		log.Println(addr)
		// Check IP dạng IPv4, không phải loopback
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no valid local IP found")
}

func main() {
	ip, err := getLocalIP()
	if err != nil {
		log.Fatalf("Error getting local IP: %v", err)
	}
	fmt.Println(ip)
}

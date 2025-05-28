package controller

import (
	"crypto/rand"
	"encoding/base64"
)

// newShortID generates a random short ID.
func newShortID() string {
	b := make([]byte, 4) // 4 bytes -> 6 chars
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// remove the markers and return the data
// useful for testing
func removeMarkers(data []byte) []byte {
	b := data[1 : len(data)-2] // remove > and \r\n
	return b
}

// boolToUint8 converts a boolean to a uint8.
func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

package openapi

import (
	_ "embed"
)

//go:embed openapi.gen.yml
var specBytes []byte

// GetSpecBytes returns the embedded OpenAPI specification as a byte slice.
func GetSpecBytes() []byte {
	return specBytes
}

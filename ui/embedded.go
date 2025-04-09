package ui

import "embed"

//go:embed dist
var dist embed.FS

func GetDist() embed.FS {
	return dist
}

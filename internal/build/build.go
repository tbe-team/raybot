package build

import "runtime"

// Build info will be linked by ldflags during build
var (
	Version   string
	Date      string
	GoVersion string
)

type Info struct {
	Version   string
	Date      string
	GoVersion string
}

func GetBuildInfo() Info {
	return Info{
		Version:   Version,
		Date:      Date,
		GoVersion: runtime.Version(),
	}
}

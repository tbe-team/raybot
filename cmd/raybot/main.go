package main

import (
	"flag"
	"os"

	"github.com/tbe-team/raybot/cmd/raybot/standalone"
)

func main() {
	var (
		configFilePath string
		dbPath         string
	)

	flag.StringVar(&configFilePath, "config", "", "path to the config file")
	flag.StringVar(&dbPath, "db", "", "path to the SQLite database file")
	flag.Parse()

	if configFilePath == "" || dbPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	standalone.Run(configFilePath, dbPath)
}

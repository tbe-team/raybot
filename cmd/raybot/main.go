package main

import (
	"log"
	"os"

	"github.com/tbe-team/raybot/cmd/raybot/pic"
	"github.com/tbe-team/raybot/internal/application"
	"github.com/tbe-team/raybot/pkg/cmdutil"
)

func main() {
	cfg, err := application.LoadConfig(os.Args[1])
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	app, cleanup, err := application.New(cfg)
	if err != nil {
		log.Printf("failed to create application: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := cleanup(); err != nil {
			log.Printf("failed to cleanup application: %v\n", err)
			os.Exit(1)
		}
	}()

	interruptChan := cmdutil.InterruptChan()

	go func() {
		if err := pic.Start(app); err != nil {
			log.Printf("failed to start PIC serial service: %v\n", err)
			os.Exit(1)
		}
	}()

	<-interruptChan
}

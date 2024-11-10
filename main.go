package main

import (
	"flag"
	"fmt"
	"os"
	"pedalboard/internal/config"
	"pedalboard/internal/server"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the config file")
	flag.Parse()

	config, err := config.NewFromFile(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	svr := server.New(config)
	if err := svr.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}

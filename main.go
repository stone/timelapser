package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	mode := flag.String("mode", "snapshot", "mode: snapshot or timelapse")
	continuous := flag.Bool("continuous", false, "run continuously with interval from config")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	switch *mode {
	case "snapshot":
		if *continuous {
			for {
				if err := takeSnapshot(config); err != nil {
					fmt.Fprintf(os.Stderr, "Error taking snapshot: %v\n", err)
				}
				time.Sleep(time.Duration(config.Interval) * time.Second)
			}
		} else {
			if err := takeSnapshot(config); err != nil {
				fmt.Fprintf(os.Stderr, "Error taking snapshot: %v\n", err)
				os.Exit(1)
			}
		}
	case "timelapse":
		if err := createTimelapse(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating timelapse: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s\n", *mode)
		os.Exit(1)
	}
}

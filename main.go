package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

var logger *slog.Logger

func main() {
	flagConfigPath := flag.String("config", "config.yaml", "path to config file")
	flagSnapshot := flag.Bool("snapshot", false, "Do a snapshot of configured cameras once")
	flagTimelapse := flag.Bool("timelapse", false, "Create timelapse for configured cameras")
	flagLogLevel := flag.String("log", "INFO", "Log level (DEBUG, INFO)")
	flag.Parse()
	var loglevel slog.Level
	switch *flagLogLevel {
	case "DEBUG":
		loglevel = slog.LevelDebug
	case "INFO":
		loglevel = slog.LevelInfo
	}

	logger = slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:   loglevel,
			NoColor: !isatty.IsTerminal(os.Stderr.Fd()),
		}),
	)

	logger.Debug("Opening configuration file", "config", *flagConfigPath)
	config, err := loadConfig(*flagConfigPath)
	if err != nil {
		logger.Error("Error loading config", "config", *flagConfigPath, "error", err)
		os.Exit(1)
	}

	if logger.Handler().Enabled(nil, slog.LevelDebug) {
		for _, camConfig := range config.Cameras {
			logger.Debug("Camera configuration", "camera", camConfig.Name, "url", camConfig.SnapshotURL, "delete", camConfig.Delete)
		}
	}

	// Take snapshot of cameras and quit
	if *flagSnapshot {
		if err := takeSnapshot(config); err != nil {
			logger.Error("Error taking snapshot", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *flagTimelapse {
		if err := createTimelapse(config); err != nil {
			logger.Error("Error creating timelapse", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Default to run forever
	for {
		if err := takeSnapshot(config); err != nil {
			logger.Error("Error taking snapshot", "error", err, "resume", true)
		}
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

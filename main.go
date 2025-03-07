package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/robfig/cron/v3"
	"github.com/stone/timelapser/internal/camera"
	"github.com/stone/timelapser/internal/config"
	"github.com/stone/timelapser/internal/snapshot"
	"github.com/stone/timelapser/internal/timelapse"
)

// Build information. Populated at build-time using -ldflags
var (
	Version   string
	GitCommit string
)

var logger *slog.Logger

func main() {
	flagConfigPath := flag.String("config", "config.yaml", "path to config file")
	flagSnapshot := flag.Bool("snapshot", false, "Do a single snapshot of all configured cameras")
	flagTimelapse := flag.Bool("timelapse", false, "Create timelapse for all configured cameras and quit, (images not deleted)")
	flagLogLevel := flag.String("log", "INFO", "Log level (DEBUG, INFO)")
	flagListCameras := flag.Bool("list", false, "List configured cameras")
	flagGetConfig := flag.Bool("example-config", false, "Print example configuration to stdout")
	flagGetVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *flagGetVersion {
		fmt.Printf("timelapser %s (%s)\n", Version, GitCommit)
		os.Exit(0)
	}

	if *flagGetConfig {
		fmt.Println(config.NewExampleConfig())
		os.Exit(0)
	}

	crn := cron.New()

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
	config, err := config.LoadConfig(*flagConfigPath, logger)
	if err != nil {
		logger.Error("Error loading config", "config", *flagConfigPath, "error", err)
		os.Exit(1)
	}

	if logger.Handler().Enabled(context.TODO(), slog.LevelDebug) {
		for _, camConfig := range config.Cameras {
			logger.Debug("Camera configuration", "camera", camConfig.Name, "url", camConfig.SnapshotURL, "delete", camConfig.Delete)
		}
	}

	if *flagListCameras {
		camera.ListCameras(config)
		os.Exit(0)
	}

	logger.Info("Starting timelapser", "version", Version, "git", GitCommit)

	// Take snapshot of cameras and quit
	if *flagSnapshot {
		if err := snapshot.TakeSnapshot(config); err != nil {
			logger.Error("Error taking snapshot", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *flagTimelapse {
		if err := timelapse.CreateAllTimelapse(config, logger); err != nil {
			logger.Error("Error creating timelapse", "error", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Schedule snapshots
	for _, camConfig := range config.Cameras {
		// Simple way of picking up default inteval or use camera specific interval
		timelapseInterval := config.TimelapseInterval
		if camConfig.TimelapseInterval != "" {
			timelapseInterval = camConfig.TimelapseInterval
		}

		interval := config.Interval
		if camConfig.Interval != "" {
			interval = camConfig.Interval
		}

		logger.Info("Scheduling camera snapshot", "name", camConfig.Name, "interval", interval)
		crn.AddFunc(interval, func() {
			if err := snapshot.TakeCameraSnapshot(&camConfig, config.OutputDir, logger); err != nil {
				logger.Error("Error taking snapshot", "name", camConfig.Name, "error", err)
			}
		})

		logger.Info("Scheduling timelapse generation", "name", camConfig.Name, "timelapseInterval", timelapseInterval)
		crn.AddFunc(timelapseInterval, func() {
			if err := timelapse.CreateTimelapse(&camConfig, config.OutputDir, logger); err != nil {
				logger.Error("Error generating timelapse", "name", camConfig.Name, "error", err)
			}
		})

	}

	// Start the scheduler
	crn.Start()

	// Keep the program running
	select {}
}

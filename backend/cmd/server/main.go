package main

import (
	"flag"
	"log/slog"
	"os"
	"stargazer/internal/core"
	"stargazer/internal/server"
)

type arguments struct {
	configPath string
}

func newArguments() *arguments {
	var configPath string
	flag.StringVar(&configPath, "config-path", "stargazer-config.toml", "specify path to configuration file")
	flag.Parse()

	return &arguments{configPath: configPath}
}

func setupLogger() {
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, &opts)
	slog.SetDefault(slog.New(handler))
}

func main() {
	args := newArguments()
	setupLogger()

	config, err := core.LoadConfig(args.configPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(111)
	}
	appState, err := core.CreateAppState(config)
	if err != nil {
		slog.Error("Failed to create app state", "error", err)
		os.Exit(111)
	}

	engine := server.CreateEngine(appState)
	engine.Run(config.Server.BindAddress())
}

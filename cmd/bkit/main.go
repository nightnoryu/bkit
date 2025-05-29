package main

import (
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/nightnoryu/bkit/internal/dockerfile"
	appconfig "github.com/nightnoryu/bkit/internal/frontend/app/config"
)

const (
	appID = "bkit"
)

// These variables come from -ldflags settings
// Here also setup their fallback values
var (
	Commit          = "UNKNOWN"
	DockerfileImage = string(dockerfile.Dockerfile14) // default image for dockerfile
)

func main() {
	ctx := context.Background()

	ctx = subscribeForKillSignals(ctx)

	err := runApp(ctx, os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(ctx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	configPath, err := appconfig.DefaultConfigPath()
	if err != nil {
		configPath = "" // Ignore err if default path unavailable
	}

	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	app := &cli.App{
		Name:  appID,
		Usage: "Container-native build system",
		Commands: []*cli.Command{
			build(workdir),
			config(),
			version(),
			cache(),
			fmtCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "bkit config",
				Aliases: []string{"c"},
				EnvVars: []string{"BKIT_CONFIG"},
				Value:   configPath,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Verbose output to stderr",
				Aliases: []string{"v"},
			},
			&cli.StringFlag{
				Name:    "docker-config",
				Usage:   "Path to docker client config",
				Aliases: []string{"dc"},
				EnvVars: []string{"BKIT_DOCKER_CONFIG"},
			},
		},
	}

	return app.RunContext(ctx, args)
}

func subscribeForKillSignals(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			signal.Stop(ch)
		case <-ch:
		}
	}()

	return ctx
}

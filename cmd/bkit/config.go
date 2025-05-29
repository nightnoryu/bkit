package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/urfave/cli/v2"

	appconfig "github.com/nightnoryu/bkit/internal/frontend/app/config"
	infraconfig "github.com/nightnoryu/bkit/internal/frontend/infrastructure/config"
)

func config() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Manipulate bkit config",
		Subcommands: []*cli.Command{
			configInit(),
		},
	}
}

func configInit() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Create default bkit config",
		Action: func(ctx *cli.Context) error {
			var opts commonOpt
			opts.scan(ctx)

			logger := makeLogger(opts.verbose)

			defaultConfig, err := infraconfig.Parser{}.Dump(appconfig.DefaultConfig)
			if err != nil {
				return err
			}

			defaultConfigBuffer := &bytes.Buffer{}
			err = json.Indent(defaultConfigBuffer, defaultConfig, "", "    ")
			if err != nil {
				return err
			}

			configPath := opts.configPath
			configDir := path.Dir(configPath)

			err = os.MkdirAll(configDir, 0o755)
			if err != nil {
				return fmt.Errorf("failed to create folder for config %s: %w", configDir, err)
			}

			err = os.WriteFile(configPath, defaultConfigBuffer.Bytes(), 0o600)
			if err != nil {
				return fmt.Errorf("failed to write file for config %s: %w", configDir, err)
			}

			logger.Outputf("Default config created in %s\n", configPath)

			return nil
		},
	}
}

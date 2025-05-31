package main

import (
	"encoding/json"

	"github.com/urfave/cli/v2"

	appversion "github.com/nightnoryu/bkit/internal/frontend/app/version"
)

func version() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "Show bkit version info",
		Action: executeVersion,
	}
}

func executeVersion(ctx *cli.Context) error {
	var opt commonOpt
	opt.scan(ctx)

	logger := makeLogger(opt.verbose)

	v := struct {
		Version    string `json:"version"`
		APIVersion string `json:"apiVersion"`
		Commit     string `json:"commit"`
		Dockerfile string `json:"dockerfile"`
	}{
		Version:    Version,
		APIVersion: appversion.APIVersionV1,
		Commit:     Commit,
		Dockerfile: DockerfileImage,
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	logger.Outputf(string(bytes) + "\n")

	return nil
}

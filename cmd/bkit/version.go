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
		Commit     string `json:"commit"`
		Date       string `json:"date"`
		BuiltBy    string `json:"built_by"`
		APIVersion string `json:"apiVersion"`
		Dockerfile string `json:"dockerfile"`
	}{
		Version:    Version,
		Commit:     Commit,
		Date:       Date,
		BuiltBy:    BuiltBy,
		APIVersion: appversion.APIVersionV1,
		Dockerfile: DockerfileImage,
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	logger.Outputf(string(bytes) + "\n")

	return nil
}

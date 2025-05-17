package jsonnet

import (
	"fmt"
	"os"
	"path"

	jsonnetformatter "github.com/google/go-jsonnet/formatter"
)

type Formatter struct{}

func (formatter Formatter) Format(configPath string) (string, error) {
	filename := path.Base(configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	options := jsonnetformatter.DefaultOptions()
	options.Indent = 4
	options.StringStyle = jsonnetformatter.StringStyleLeave

	data, err := jsonnetformatter.Format(filename, string(configData), options)
	if err != nil {
		return "", fmt.Errorf("failed to format config file: %w", err)
	}

	return data, nil
}

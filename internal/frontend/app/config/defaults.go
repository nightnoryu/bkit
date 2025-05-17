package config

import (
	"fmt"
	"os"
)

const (
	defaultConfigPath = "%s/.brewkit/config"
)

var (
	DefaultConfig = Config{}
)

func DefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to receive user home dir: %w", err)
	}

	return fmt.Sprintf(defaultConfigPath, homeDir), nil
}

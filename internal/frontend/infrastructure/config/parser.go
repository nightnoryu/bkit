package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/google/go-jsonnet"

	"github.com/nightnoryu/bkit/internal/common/slices"
	"github.com/nightnoryu/bkit/internal/frontend/app/config"
)

type Parser struct{}

func (p Parser) Config(configPath string) (config.Config, error) {
	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config.Config{}, config.ErrConfigNotFound
		}

		return config.Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	vm := jsonnet.MakeVM()

	jsonnet.Version()

	data, err := vm.EvaluateAnonymousSnippet(path.Base(configPath), string(fileBytes))
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to compile jsonnet for config: %w", err)
	}

	var c Config

	err = json.Unmarshal([]byte(data), &c)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse json config: %w", err)
	}

	return config.Config{
		Secrets: slices.Map(c.Secrets, func(s Secret) config.Secret {
			return config.Secret{
				ID:   s.ID,
				Path: os.ExpandEnv(s.Path),
			}
		}),
	}, nil
}

func (p Parser) Dump(srcConfig config.Config) ([]byte, error) {
	c := Config{
		Secrets: slices.Map(srcConfig.Secrets, func(s config.Secret) Secret {
			return Secret{
				ID:   s.ID,
				Path: s.Path,
			}
		}),
	}

	data, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config to json: %w", err)
	}

	return data, nil
}

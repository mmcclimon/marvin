package marvin

import (
	"errors"

	"github.com/BurntSushi/toml"
)

var ErrShuttingDown = errors.New("shutting down")

func FromFile(path string, registry Registry) (*Hub, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)

	if err != nil {
		return nil, err
	}

	return cfg.Assemble(registry)
}

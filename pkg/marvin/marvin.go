package marvin

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Marvin struct {
	err error
	cfg Config
}

type Config struct {
	Name string
}

func FromFile(path string) *Marvin {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)

	return &Marvin{
		cfg: cfg,
		err: err,
	}
}

func (m *Marvin) Run() error {
	if m.err != nil {
		return m.err
	}

	fmt.Printf("%+v\n", m)
	return nil
}

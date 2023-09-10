package marvin

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Marvin struct {
	err      error
	buses    map[string]Bus
	reactors map[string]Reactor
}

func FromFile(path string) *Marvin {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)

	if err != nil {
		return &Marvin{err: err}
	}

	return cfg.Assemble()
}

func (m *Marvin) Run() error {
	if m.err != nil {
		return m.err
	}

	fmt.Printf("%+v\n", m)
	return nil
}

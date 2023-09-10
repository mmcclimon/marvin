package registry

import (
	"fmt"

	"github.com/mmcclimon/marvin"
)

type Registry struct {
	buses    map[string]marvin.BusAssembler
	reactors map[string]marvin.ReactorAssembler
}

var singleton = Registry{
	buses:    make(map[string]marvin.BusAssembler),
	reactors: make(map[string]marvin.ReactorAssembler),
}

func Default() Registry { return singleton }

func (r Registry) hasBus(name string) bool {
	_, ok := r.buses[name]
	return ok
}

func (r Registry) hasReactor(name string) bool {
	_, ok := r.reactors[name]
	return ok
}

func (r Registry) ReactorFor(name string) marvin.ReactorAssembler {
	return singleton.reactors[name]
}

func (r Registry) BusFor(name string) marvin.BusAssembler {
	return singleton.buses[name]
}

func RegisterReactor(name string, assembler marvin.ReactorAssembler) {
	if singleton.hasReactor(name) {
		panic(fmt.Sprintf("cannot register duplicate reactor '%s'", name))
	}

	singleton.reactors[name] = assembler
}

func RegisterBus(name string, assembler marvin.BusAssembler) {
	if singleton.hasBus(name) {
		panic(fmt.Sprintf("cannot register duplicate bus '%s'", name))
	}

	singleton.buses[name] = assembler
}

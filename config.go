package marvin

import (
	"fmt"
	"strings"
)

type arbitraryConfig = map[string]any

type Config struct {
	Name    string
	Bus     map[string]arbitraryConfig
	Reactor map[string]arbitraryConfig
	err     assemblyError
}

type assemblyError struct {
	errs []error
}

func (cfg *Config) Assemble() *Marvin {
	// errMarv := func(err error) *Marvin { return &Marvin{err: err} }
	marv := &Marvin{
		events:   make(chan Event),
		errs:     make(chan error),
		reactors: make(map[string]Reactor),
		buses:    make(map[string]Bus),
	}

	cfg.assembleBuses(marv)
	cfg.assembleReactors(marv)

	if cfg.err.hasErrors() {
		marv.err = cfg.err
	}

	return marv
}

func (cfg *Config) assembleBuses(marv *Marvin) {
	for name, busConfig := range cfg.Bus {
		assembler, err := extractAssembler("bus", name, busConfig, busRegistry)
		if err != nil {
			cfg.err.add(err)
			continue
		}

		bus, err := assembler(busConfig)
		if err != nil {
			cfg.err.add(fmt.Errorf("error assembling bus '%s': %w", name, err))
			continue
		}

		marv.buses[name] = bus
	}
}

func (cfg *Config) assembleReactors(marv *Marvin) {
	for name, reactorConfig := range cfg.Reactor {
		assembler, err := extractAssembler("reactor", name, reactorConfig, reactorRegistry)
		if err != nil {
			cfg.err.add(err)
			continue
		}

		reactor, err := assembler(reactorConfig)
		if err != nil {
			cfg.err.add(fmt.Errorf("error assembling reactor '%s': %w", name, err))
			continue
		}

		marv.reactors[name] = reactor
	}
}

type componentAssembler interface {
	BusAssembler | ReactorAssembler
}

func extractAssembler[T componentAssembler](
	thing string,
	name string,
	conf arbitraryConfig,
	registry map[string]T,
) (T, error) {
	typ, ok := conf["type"]
	if !ok {
		return nil, fmt.Errorf("no type found for %s '%s'", thing)
	}

	typString, ok := typ.(string)
	if !ok {
		return nil, fmt.Errorf("'type' for %s '%s' is not a string", thing, name)
	}

	assembler, ok := registry[typString]
	if !ok {
		return nil, fmt.Errorf("unknown type found for %s '%s'", thing, name)
	}

	return assembler, nil

}

func (ae *assemblyError) add(err error) {
	ae.errs = append(ae.errs, err)
}

func (ae *assemblyError) hasErrors() bool {
	return len(ae.errs) > 0
}

func (ae assemblyError) Error() string {
	all := make([]string, len(ae.errs))

	for i := 0; i < len(all); i++ {
		all[i] = "    " + ae.errs[i].Error()
	}

	return "assembly errors:\n" + strings.Join(all, "\n")
}

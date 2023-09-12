package marvin

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type arbitraryConfig = map[string]any

type Config struct {
	Name     string
	LogLevel slog.Level `toml:"log_level"`
	Bus      map[string]arbitraryConfig
	Reactor  map[string]arbitraryConfig
	err      assemblyError
}

type Registry interface {
	BusFor(string) BusAssembler
	ReactorFor(string) ReactorAssembler
}

type assemblyError struct {
	errs []error
}

func (cfg *Config) Assemble(registry Registry) *Marvin {
	marv := &Marvin{
		events:   make(chan Event),
		errs:     make(chan error),
		reactors: make(map[ReactorName]Reactor),
		buses:    make(map[BusName]Bus),
	}

	cfg.assembleBuses(marv, registry)
	cfg.assembleReactors(marv, registry)

	if cfg.err.hasErrors() {
		marv.err = cfg.err
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	slog.SetDefault(logger)

	return marv
}

func (cfg *Config) assembleBuses(marv *Marvin, registry Registry) {
	for name, busConfig := range cfg.Bus {
		assembler, err := extractAssembler("bus", name, busConfig, registry.BusFor)
		if err != nil {
			cfg.err.add(err)
			continue
		}

		identifier := BusName(name)
		bus, err := assembler(identifier, busConfig)
		if err != nil {
			cfg.err.add(fmt.Errorf("error assembling bus '%s': %w", name, err))
			continue
		}

		marv.buses[identifier] = bus
	}
}

func (cfg *Config) assembleReactors(marv *Marvin, registry Registry) {
	for name, reactorConfig := range cfg.Reactor {
		assembler, err := extractAssembler("reactor", name, reactorConfig, registry.ReactorFor)
		if err != nil {
			cfg.err.add(err)
			continue
		}

		identifier := ReactorName(name)
		reactor, err := assembler(identifier, reactorConfig)
		if err != nil {
			cfg.err.add(fmt.Errorf("error assembling reactor '%s': %w", name, err))
			continue
		}

		marv.reactors[identifier] = reactor
	}
}

type componentAssembler interface {
	BusAssembler | ReactorAssembler
}

func extractAssembler[T componentAssembler](
	ct string,
	name string,
	rawConf arbitraryConfig,
	fetcher func(string) T,
) (T, error) {
	var conf struct{ Type string }

	if err := mapstructure.Decode(rawConf, &conf); err != nil {
		return nil, fmt.Errorf("could not extract type for %s '%s': %w", ct, name)
	}

	delete(rawConf, "type")

	assembler := fetcher(conf.Type)
	if assembler == nil {
		return nil, fmt.Errorf("unknown type found for %s '%s'", ct, name)
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

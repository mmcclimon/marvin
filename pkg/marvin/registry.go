package marvin

import "fmt"

var (
	busRegistry     = make(map[string]BusAssembler)
	reactorRegistry = make(map[string]ReactorAssembler)
)

func RegisterReactor(name string, assembler ReactorAssembler) {
	_, exists := reactorRegistry[name]
	if exists {
		panic(fmt.Sprintf("cannot register duplicate reactor '%s'", name))
	}

	reactorRegistry[name] = assembler
}

func RegisterBus(name string, assembler BusAssembler) {
	_, exists := busRegistry[name]
	if exists {
		panic(fmt.Sprintf("cannot register duplicate bus '%s'", name))
	}

	busRegistry[name] = assembler
}

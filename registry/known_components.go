package registry

import (
	"github.com/mmcclimon/marvin/buses/term"
	"github.com/mmcclimon/marvin/reactors/echo"
)

// RegisterAllKnownComponents adds all the default buses and reactors with
// their well-known names (i.e., buses/term gets registered as "term",
// reactors/echo as "echo", and so on).
func RegisterAllKnownComponents() {
	RegisterBus("term", term.Assemble)

	RegisterReactor("echo", echo.Assemble)
}

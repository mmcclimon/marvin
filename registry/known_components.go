package registry

import (
	"github.com/mmcclimon/marvin/buses/discord"
	"github.com/mmcclimon/marvin/buses/term"
	"github.com/mmcclimon/marvin/reactors/echo"
	"github.com/mmcclimon/marvin/reactors/eject"
	"github.com/mmcclimon/marvin/reactors/uptime"
)

// RegisterAllKnownComponents adds all the default buses and reactors with
// their well-known names (i.e., buses/term gets registered as "term",
// reactors/echo as "echo", and so on).
func RegisterAllKnownComponents() {
	RegisterBus("term", term.Assemble)
	RegisterBus("discord", discord.Assemble)

	RegisterReactor("echo", echo.Assemble)
	RegisterReactor("eject", eject.Assemble)
	RegisterReactor("uptime", uptime.Assemble)
}

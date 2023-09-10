package echo

import "github.com/mmcclimon/marvin/pkg/marvin"

type Echo struct{}

func Assemble(cfg any) (marvin.Reactor, error) {
	return &Echo{}, nil
}

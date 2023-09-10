package term

import "github.com/mmcclimon/marvin/pkg/marvin"

type Term struct{}

func Assemble(cfg any) (marvin.Bus, error) {
	return &Term{}, nil
}

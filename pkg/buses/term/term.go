package term

import (
	"context"

	"github.com/mmcclimon/marvin/pkg/marvin"
)

type Term struct{}

func Assemble(cfg any) (marvin.Bus, error) {
	return &Term{}, nil
}

func (term *Term) Run(ctx context.Context) error {
	return nil
}

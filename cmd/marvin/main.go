package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/mmcclimon/marvin"
	"github.com/mmcclimon/marvin/registry"
)

var configFlag = flag.String("c", "", "path to config file")

func main() {
	registry.RegisterAllKnownComponents()

	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	err := marvin.FromFile(*configFlag, registry.Default()).Run()
	if errors.Is(err, marvin.ErrShuttingDown) {
		slog.Info("bye now!")
	} else if err != nil {
		fmt.Println("fatal error:")
		fmt.Println(err)
		os.Exit(1)
	}
}

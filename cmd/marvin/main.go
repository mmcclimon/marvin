package main

import (
	"errors"
	"flag"
	"log"
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
		log.Println("bye now!")
	} else if err != nil {
		log.Fatal(err)
	}
}

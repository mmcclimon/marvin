package main

import (
	"flag"
	"log"
	"os"

	"github.com/mmcclimon/marvin/pkg/buses/term"
	"github.com/mmcclimon/marvin/pkg/marvin"

	"github.com/mmcclimon/marvin/pkg/reactors/echo"
)

var configFlag = flag.String("c", "", "path to config file")

func main() {
	registerComponents()

	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	err := marvin.FromFile(*configFlag).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func registerComponents() {
	marvin.RegisterBus("term", term.Assemble)

	marvin.RegisterReactor("echo", echo.Assemble)
}

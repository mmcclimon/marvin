package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/mmcclimon/marvin"
	"github.com/mmcclimon/marvin/buses/term"
	"github.com/mmcclimon/marvin/reactors/echo"
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
	if errors.Is(err, marvin.ErrShuttingDown) {
		log.Println("bye now!")
	} else if err != nil {
		log.Fatal(err)
	}
}

func registerComponents() {
	marvin.RegisterBus("term", term.Assemble)
	marvin.RegisterReactor("echo", echo.Assemble)
}

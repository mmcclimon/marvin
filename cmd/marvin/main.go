package main

import (
	"flag"
	"log"
	"os"

	"github.com/mmcclimon/marvin/pkg/marvin"
)

var configFlag = flag.String("c", "", "path to config file")

func main() {
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

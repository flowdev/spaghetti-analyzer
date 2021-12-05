package main

import (
	"log"
	"os"

	"github.com/flowdev/spaghetti-analyzer/deps/testdata/half-pkgs-proj/db/store"
	"github.com/flowdev/spaghetti-analyzer/deps/testdata/half-pkgs-proj/domain1"
	"github.com/flowdev/spaghetti-analyzer/deps/testdata/half-pkgs-proj/domain2"
)

func main() {
	doIt(os.Args[1:])
}

func doIt(args []string) {
	log.Printf("INFO - this is the main package, args: %q", args)
	s := &store.Store{}
	domain1.HandleDomain1Route1(s)
	domain1.HandleDomain1Route2(s)

	domain2.HandleDomain2Route1(s)
	domain2.HandleDomain2Route2(s)
}

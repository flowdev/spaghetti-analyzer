package main

import (
	"strconv"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestAnalyze(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"analyze": func(ts *testscript.TestScript, _ bool, args []string) {
				expectedReturnCode, err := strconv.Atoi(args[0])
				if err != nil {
					ts.Fatalf("fatal return code error (%q): %v", args[0], err)
				}

				args = args[1:]
				actualReturnCode := analyze(args)

				if actualReturnCode != expectedReturnCode {
					ts.Fatalf("Expected return code %d but got: %d", expectedReturnCode, actualReturnCode)
				}
			},
		},
		TestWork: false,
	})
}

package tree_test

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/flowdev/spaghetti-analyzer/parse"
	"github.com/flowdev/spaghetti-analyzer/tree"
)

func TestGenerate(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"dirTree": callGenerate,
		},
		// TestWork: true,
	})
}

func callGenerate(ts *testscript.TestScript, _ bool, args []string) {
	workDir := ts.Getenv("WORK")
	treeFile := workDir + "/dirtree.actual"
	packs, err := parse.DirTree(workDir)
	if err != nil {
		ts.Fatalf("ERROR: %v", err)
	}

	name := "project-root"
	if len(args) > 0 {
		name = args[0]
	}

	output, err := tree.Generate(workDir, name, packs)
	if err != nil {
		ts.Fatalf("ERROR: Unable to generate directory tree: %v", err)
	}

	err = os.WriteFile(treeFile, []byte(output+"\n"), 0666)
	if err != nil {
		ts.Fatalf("ERROR: Unable to write file '%s': %v", treeFile, err)
	}
}

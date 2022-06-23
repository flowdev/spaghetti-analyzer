package doc_test

import (
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-analyzer/doc"
	"github.com/flowdev/spaghetti-cutter/data"
)

var testDepMap = analdata.DependencyMap{
	"a": analdata.PkgImports{
		PkgType: data.TypeGod,
		Imports: map[string]data.PkgType{
			"b/c/d":   data.TypeTool,
			"epsilon": data.TypeTool,
			"escher":  data.TypeTool,
			"f":       data.TypeStandard,
			"z":       data.TypeDB,
		},
	},
	"z": analdata.PkgImports{
		PkgType: data.TypeDB,
		Imports: map[string]data.PkgType{
			"b/c/d":   data.TypeTool,
			"epsilon": data.TypeTool,
			"escher":  data.TypeTool,
			"x":       data.TypeDB,
		},
	},
	"x": analdata.PkgImports{
		PkgType: data.TypeDB,
		Imports: map[string]data.PkgType{
			"b/c/d":  data.TypeTool,
			"escher": data.TypeTool,
		},
	},
	"m": analdata.PkgImports{
		PkgType: data.TypeStandard,
		Imports: map[string]data.PkgType{
			"b/c/d": data.TypeTool,
			"x":     data.TypeDB,
		},
	},
	"f": analdata.PkgImports{
		PkgType: data.TypeStandard,
		Imports: map[string]data.PkgType{
			"f/g": data.TypeStandard,
			"f/h": data.TypeStandard,
			"f/i": data.TypeStandard,
		},
	},
	"f/g": analdata.PkgImports{
		PkgType: data.TypeStandard,
		Imports: map[string]data.PkgType{
			"escher": data.TypeTool,
			"f/j":    data.TypeStandard,
			"x":      data.TypeDB,
		},
	},
	"f/h": analdata.PkgImports{
		PkgType: data.TypeStandard,
		Imports: map[string]data.PkgType{
			"m": data.TypeStandard,
			"x": data.TypeDB,
		},
	},
	"f/i": analdata.PkgImports{
		PkgType: data.TypeStandard,
		Imports: map[string]data.PkgType{
			"escher": data.TypeTool,
		},
	},
}

func TestWriteDocs(t *testing.T) {
	rootPkg := "github.com/flowdev/tst"
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"writeDocs": func(ts *testscript.TestScript, _ bool, args []string) {
				workDir := ts.Getenv("WORK")

				if len(args) != 1 {
					ts.Fatalf("ERROR: Expected 1 argument (givenDependencyTablePkgs) but got: : %q", args)
				}
				givenDependencyTablePkgs := strings.Split(args[0], ",")

				doc.WriteDocs(givenDependencyTablePkgs, testDepMap, rootPkg, workDir)
			},
		},
		TestWork: false,
	})
}

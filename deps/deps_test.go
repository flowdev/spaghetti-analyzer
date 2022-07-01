package deps_test

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-analyzer/deps"
	"github.com/flowdev/spaghetti-analyzer/parse"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
	"github.com/flowdev/spaghetti-cutter/config"
	"github.com/flowdev/spaghetti-cutter/data"
)

func TestFill(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"fill": func(ts *testscript.TestScript, _ bool, args []string) {
				workDir := ts.Getenv("WORK")
				depsFile := workDir + "/deps.actual"

				if len(args) != 1 {
					ts.Fatalf("ERROR: Expected 1 argument (givenConfig) but got: : %q", args)
				}
				givenConfig := args[0]
				cfg, err := config.Parse([]byte(givenConfig), ".spaghetti-cutter.hjson")
				if err != nil {
					ts.Fatalf("fatal config error: %v", err)
				}

				packs, err := parse.DirTree(workDir)
				if err != nil {
					ts.Fatalf("fatal parse error: %v", err)
				}

				depMap := make(analdata.DependencyMap, 256)
				rootPkg := parse.RootPkg(packs)
				ts.Logf("root package: %s", rootPkg)
				pkgInfos := pkgs.UniquePackages(packs)
				for _, pkgInfo := range pkgInfos {
					deps.Fill(pkgInfo.Pkg, rootPkg, cfg, &depMap)
				}

				sDeps := prettyPrint(depMap)
				err = os.WriteFile(depsFile, []byte(sDeps+"\n"), 0666)
				if err != nil {
					ts.Fatalf("ERROR: Unable to write file '%s': %v", depsFile, err)
				}
			},
		},
		TestWork: false,
	})
}

func prettyPrint(deps analdata.DependencyMap) string {
	sb := strings.Builder{}

	for _, pkg := range deps.SortedPkgNames() {
		imps := deps[pkg]
		pkgTypeRune := data.TypeLetter(imps.PkgType)

		for _, imp := range sortedImpNames(imps.Imports) {
			sb.WriteString(pkg)
			sb.WriteString(" [")
			sb.WriteRune(pkgTypeRune)
			sb.WriteString("] imports: ")
			sb.WriteString(imp)
			sb.WriteString(" [")
			sb.WriteRune(data.TypeLetter(imps.Imports[imp]))
			sb.WriteString("]\n")
		}
	}

	return sb.String()
}

func sortedImpNames(imps map[string]data.PkgType) []string {
	names := make([]string, 0, len(imps))
	for imp := range imps {
		names = append(names, imp)
	}
	sort.Strings(names)
	return names
}

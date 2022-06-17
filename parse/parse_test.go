package parse_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/parse"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestDirTree(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"parseDirTree": callParseDirTree,
		},
		TestWork: true,
	})
}

func callParseDirTree(ts *testscript.TestScript, _ bool, args []string) {
	workDir := ts.Getenv("WORK")
	pkgFile := workDir + "/packages.actual"
	output := ""

	actualPkgs, err := parse.DirTree(workDir)
	if err != nil {
		ts.Logf("ERROR: %v", err)
		output = "error: true"
	} else {
		output = "error: false\n" + packagesAsString(actualPkgs)
	}

	err = os.WriteFile(pkgFile, []byte(output+"\n\n"), 0666)
	if err != nil {
		ts.Fatalf("ERROR: Unable to write file '%s': %v", pkgFile, err)
	}
}
func packagesAsString(packs []*pkgs.Package) string {
	strPkgs := make([]string, len(packs))

	for i, p := range packs {
		strPkgs[i] = p.Name + ": " + p.PkgPath
		if isTestPkg(p) {
			strPkgs[i] += " [T]"
		}
	}
	sort.Strings(strPkgs)
	return strings.Join(strPkgs, "\n")
}

func TestRootPkg(t *testing.T) {
	specs := []struct {
		name          string
		givenPkgPaths []string
		expectedRoot  string
	}{
		{
			name:          "empty",
			givenPkgPaths: []string{"", ""},
			expectedRoot:  "",
		}, {
			name:          "nothing-in-common",
			givenPkgPaths: []string{"a", "ba"},
			expectedRoot:  "",
		}, {
			name:          "test-packages",
			givenPkgPaths: []string{"pkg/x/a", "pkg/x/a_test", "pkg/x/a.test"},
			expectedRoot:  "pkg/x/a",
		}, {
			name:          "x-packages",
			givenPkgPaths: []string{"pkg/x/a", "pkg/x/b", "pkg/x/c"},
			expectedRoot:  "pkg/x/",
		}, {
			name:          "all-on-github",
			givenPkgPaths: []string{"github.com/org1/proj1", "github.com/org1/proj2", "github.com/borg2/proj3"},
			expectedRoot:  "github.com/",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualRoot := parse.RootPkg(pkgsForPaths(spec.givenPkgPaths))
			if actualRoot != spec.expectedRoot {
				t.Errorf("expected common root %q, got %q",
					spec.expectedRoot, actualRoot)
			}
		})
	}
}

func pkgsForPaths(paths []string) []*pkgs.Package {
	packs := make([]*pkgs.Package, len(paths))
	for i, path := range paths {
		packs[i] = &pkgs.Package{PkgPath: path}
	}
	return packs
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}

func isTestPkg(pkg *pkgs.Package) bool {
	return strings.HasSuffix(pkg.PkgPath, "_test") ||
		strings.HasSuffix(pkg.PkgPath, ".test") ||
		strings.HasSuffix(pkg.ID, ".test]") ||
		strings.HasSuffix(pkg.ID, ".test")
}

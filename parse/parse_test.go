package parse_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/parse"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
)

func TestDirTree(t *testing.T) {
	specs := []struct {
		name          string
		expectedError bool
		expectedPkgs  string
	}{
		{
			name:          "happy-path",
			expectedError: false,
			expectedPkgs: "alltst: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/alltst | " +
				"alltst: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/alltst [T] | " +
				"alltst_test: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/alltst_test [T] | " +
				"apitst: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/apitst | " +
				"apitst_test: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/apitst_test [T] | " +
				"main: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path | " +
				"main: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/alltst.test [T] | " +
				"main: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/apitst.test [T] | " +
				"main: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/unittst.test [T] | " +
				"unittst: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/unittst | " +
				"unittst: github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/unittst [T]",
		}, {
			name:          "error-path",
			expectedError: true,
			expectedPkgs:  "",
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualPkgs, err := parse.DirTree(mustAbs(filepath.Join("testdata", spec.name)))
			//t.Logf("err: %v, actualPkgs: %#v", err, actualPkgs)
			if spec.expectedError {
				if err != nil {
					t.Logf("received expected error: %v", err)
				} else {
					t.Error("expected to receive error but didn't get one")
				}
			} else if err != nil {
				t.Fatalf("received UNexpected error: %v", err)
			}
			actualPkgsString := packagesAsString(actualPkgs)
			if actualPkgsString != spec.expectedPkgs {
				t.Errorf("expected parsed packages %q, actual %q (len=%d)",
					spec.expectedPkgs, actualPkgsString, len(actualPkgs))
			}
		})
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
	return strings.Join(strPkgs, " | ")
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
				t.Errorf("expected common root %q, actual %q",
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

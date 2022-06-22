package dirs_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/doc"
	"github.com/flowdev/spaghetti-analyzer/x/dirs"
	"github.com/rogpeppe/go-internal/testscript"
)

const testFile = ".test-file"

func TestValidateRoot(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/validate-root",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"validateRoot": func(ts *testscript.TestScript, _ bool, args []string) {
				workDir := ts.Getenv("WORK")

				if len(args) != 3 {
					ts.Fatalf("ERROR: Expected 3 arguments (givenCWD, givenRoot and expectedRoot) but got: : %q", args)
				}
				givenCWD, givenRoot, expectedRoot := args[0], args[1], args[2]

				actualRoot, err := dirs.ValidateRoot(filepath.Join(workDir, givenCWD, givenRoot), testFile)
				if err != nil {
					ts.Fatalf("expected no error but got: %v", err)
				}
				if actualRoot != expectedRoot {
					ts.Fatalf("expected project root %q, got %q",
						expectedRoot, actualRoot)
				}
			},
		},
		// TestWork: true,
	})
}

func TestFindDepTables(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/find-dep-tables",
		Cmds: map[string]func(*testscript.TestScript, bool, []string){
			"findDepTables": func(ts *testscript.TestScript, _ bool, args []string) {
				workDir := ts.Getenv("WORK")
				if len(args) != 2 {
					ts.Fatalf("ERROR: Expected 2 arguments (givenStartPkgs and expectedPkgs) but got: : %q", args)
				}
				givenStartPkgs, expectedPkgs := strings.Split(args[0], ","), strings.Split(args[1], ",")
				if args[0] == "" {
					givenStartPkgs = nil
				}

				rootPkg := "github.com/flowdev/spaghetti-analyzer"
				actualPkgsWithDepTables := dirs.FindDepTables(doc.FileName, doc.Title, givenStartPkgs, workDir, rootPkg)
				ts.Logf("actualPkgsWithDepTables: %#v", actualPkgsWithDepTables)
				checkPackages(ts, expectedPkgs, actualPkgsWithDepTables)
			},
		},
		// TestWork: true,
	})
}

func checkPackages(ts *testscript.TestScript, expectedPkgs []string, actualPkgMap map[string]struct{}) {
	if len(expectedPkgs) != len(actualPkgMap) {
		ts.Fatalf("expected %d packages with dependency tables, got: %d", len(expectedPkgs), len(actualPkgMap))
	}
	packs := make([]string, 0, len(actualPkgMap))
	for p := range actualPkgMap {
		packs = append(packs, p)
	}
	sort.Strings(packs)

	for i, p := range expectedPkgs {
		if p != packs[i] {
			ts.Fatalf("expected package with dependency tables at index %d is %q, got: %q", i, p, packs[i])
		}
	}
}

func mustChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		panic(err.Error())
	}
}

func TestIncludeFile(t *testing.T) {
	specs := []struct {
		name            string
		givenName       string
		expectedInclude bool
	}{
		{
			name:            "empty",
			givenName:       "",
			expectedInclude: true,
		}, {
			name:            "normal",
			givenName:       "normal",
			expectedInclude: true,
		}, {
			name:            "vendor",
			givenName:       "vendor",
			expectedInclude: false,
		}, {
			name:            "testdata",
			givenName:       "testdata",
			expectedInclude: false,
		}, {
			name:            "dotFile",
			givenName:       ".git",
			expectedInclude: false,
		},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualInclude := dirs.IncludeFile(spec.givenName)
			if actualInclude != spec.expectedInclude {
				t.Errorf("expected include %t, got %t", spec.expectedInclude, actualInclude)
			}
		})
	}
}

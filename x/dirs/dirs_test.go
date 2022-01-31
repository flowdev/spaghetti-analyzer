package dirs_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/doc"
	"github.com/flowdev/spaghetti-analyzer/x/dirs"
)

const testFile = ".test-file"

func TestFindRoot(t *testing.T) {
	testDataDir := mustAbs(filepath.Join("testdata", "find-root"))
	specs := []struct {
		name              string
		givenCWD          string
		givenStartDir     string
		givenIgnoreVendor bool
		expectedRoot      string
	}{
		{
			name:              "given-root",
			givenCWD:          "",
			givenStartDir:     filepath.Join("in", "some", "subdir"),
			givenIgnoreVendor: false,
			expectedRoot:      filepath.Join(testDataDir, "given-root", "in"),
		}, {
			name:              "config-file",
			givenCWD:          filepath.Join("in", "some", "subdir"),
			givenStartDir:     "",
			givenIgnoreVendor: false,
			expectedRoot:      filepath.Join(testDataDir, "config-file"),
		},
	}

	initDir := mustAbs(".")
	t.Cleanup(func() {
		mustChdir(initDir)
	})
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			mustChdir(filepath.Join(testDataDir, spec.name, spec.givenCWD))

			actualRoot, err := dirs.FindRoot(spec.givenStartDir, testFile)
			if err != nil {
				t.Fatalf("expected no error but got: %v", err)
			}
			if actualRoot != spec.expectedRoot {
				t.Errorf("expected project root %q, actual %q",
					spec.expectedRoot, actualRoot)
			}
		})
	}
}

func TestFindDepTables(t *testing.T) {
	specs := []struct {
		name                      string
		givenStartPkgs            []string
		expectedPkgsWithDepTables []string
	}{
		{
			name:           "minimal",
			givenStartPkgs: nil,
			expectedPkgsWithDepTables: []string{
				"minimal",
			},
		}, {
			name:           "with-start-pkgs",
			givenStartPkgs: []string{"start-pkg"},
			expectedPkgsWithDepTables: []string{
				"start-pkg",
				"with-start-pkgs",
			},
		}, {
			name:           "with-pattern",
			givenStartPkgs: nil,
			expectedPkgsWithDepTables: []string{
				"with-pattern/bla*blue/**",
			},
		}, {
			name:           "with-root",
			givenStartPkgs: nil,
			expectedPkgsWithDepTables: []string{
				"with-root",
			},
			/*
				}, {
					name:           "many-packages",
					givenStartPkgs: nil,
					expectedPkgsWithDepTables: []string{
						"many-packages",
						"package1",
						"package2",
					},
				}, {
					name: "all-complexity",
					givenStartPkgs: []string{
						"start-pkg1", "github.com/flowdev/spaghetti-cutter/cut-it",
					},
					expectedPkgsWithDepTables: []string{
						"github.com/flowdev/spaghetti-cutter/cut-it",
						"start-pkg1",
						"all-complexity",
						"package1",
						"package2",
						"package3",
					},
			*/
		},
	}

	testDataDir := filepath.Join("testdata", "find-dep-tables")
	rootPkg := "github.com/flowdev/spaghetti-analyzer"
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualPkgsWithDepTables := dirs.FindDepTables(doc.FileName, doc.Title, spec.givenStartPkgs, filepath.Join(testDataDir, spec.name), rootPkg)
			checkPackages(t, spec.expectedPkgsWithDepTables, actualPkgsWithDepTables)
		})
	}
}

func checkPackages(t *testing.T, expectedPkgs []string, actualPkgMap map[string]struct{}) {
	if len(expectedPkgs) != len(actualPkgMap) {
		t.Errorf("expected %d packages with dependency tables, got: %d", len(expectedPkgs), len(actualPkgMap))
	}
	packs := make([]string, 0, len(actualPkgMap))
	for p := range actualPkgMap {
		packs = append(packs, p)
	}
	sort.Strings(packs)

	for i, p := range expectedPkgs {
		if p != packs[i] {
			t.Errorf("expected package with dependency tables at index %d is %q, got: %q", i, p, packs[i])
		}
	}
}

func mustChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		panic(err.Error())
	}
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}

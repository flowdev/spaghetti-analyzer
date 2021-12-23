package deps_test

import (
	"path/filepath"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/data"
	"github.com/flowdev/spaghetti-analyzer/deps"
	"github.com/flowdev/spaghetti-analyzer/parse"
	"github.com/flowdev/spaghetti-analyzer/x/config"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
)

func TestFill(t *testing.T) {
	specs := []struct {
		name         string
		givenRoot    string
		givenConfig  string
		expectedDeps string
	}{
		{
			name:        "no-config-one-pkg",
			givenRoot:   "one-pkg",
			givenConfig: `{}`,
		}, {
			name:        "no-config-only-tools",
			givenRoot:   "only-tools",
			givenConfig: `{}`,
		}, {
			name:        "no-config-standard-proj",
			givenRoot:   "standard-proj",
			givenConfig: `{}`,
		}, {
			name:        "standard-config-standard-proj",
			givenRoot:   "standard-proj",
			givenConfig: `{"tool": ["x/*"], "db": ["db/*"]}`,
		}, {
			name:        "no-config-complex-proj",
			givenRoot:   "complex-proj",
			givenConfig: `{}`,
		}, {
			name:        "standard-config-complex-proj",
			givenRoot:   "complex-proj",
			givenConfig: `{"tool": ["pkg/x/*"], "db": ["pkg/db/*"]}`,
		}, {
			name:        "no-config-half-pkgs-proj",
			givenRoot:   "half-pkgs-proj",
			givenConfig: `{}`,
		}, {
			name:        "standard-config-half-pkgs-proj",
			givenRoot:   "half-pkgs-proj",
			givenConfig: `{tool: ["x/*"], db: ["db/*"]}`,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			cfg, err := config.Parse([]byte(spec.givenConfig), spec.name)
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}

			packs, err := parse.DirTree(mustAbs(filepath.Join("testdata", spec.givenRoot)))
			if err != nil {
				t.Fatalf("Fatal parse error: %v", err)
			}

			depMap := make(data.DependencyMap, 256)
			rootPkg := parse.RootPkg(packs)
			t.Logf("root package: %s", rootPkg)
			pkgInfos := pkgs.UniquePackages(packs)
			for _, pkgInfo := range pkgInfos {
				deps.Fill(pkgInfo.Pkg, rootPkg, cfg, &depMap)
			}
			sDeps := prettyPrint(depMap)
			if sDeps != spec.expectedDeps {
				failWithDiff(t, spec.expectedDeps, sDeps)
			}
		})
	}
}

func mustAbs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err.Error())
	}
	return absPath
}

func prettyPrint(deps data.DependencyMap) string {
	return ""
}

func failWithDiff(t *testing.T, expected, actual string) {
}

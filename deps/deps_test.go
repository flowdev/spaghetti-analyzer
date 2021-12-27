package deps_test

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/deps"
	"github.com/flowdev/spaghetti-analyzer/parse"
	"github.com/flowdev/spaghetti-analyzer/x/pkgs"
	"github.com/flowdev/spaghetti-cutter/data"
	"github.com/flowdev/spaghetti-cutter/x/config"
)

func TestFill(t *testing.T) {
	specs := []struct {
		name         string
		givenRoot    string
		givenConfig  string
		expectedDeps string
	}{
		{
			name:         "no-config-one-pkg",
			givenRoot:    "one-pkg",
			givenConfig:  `{}`,
			expectedDeps: ``,
		}, {
			name:        "no-config-only-tools",
			givenRoot:   "only-tools",
			givenConfig: `{}`,
			expectedDeps: `/ [G] imports: x/tool [S]
/ [G] imports: x/tool2 [S]
x/tool2 [S] imports: x/tool [S]`,
		}, {
			name:        "standard-config-only-tools",
			givenRoot:   "only-tools",
			givenConfig: `{"tool": ["x/*"]}`,
			expectedDeps: `/ [G] imports: x/tool [T]
/ [G] imports: x/tool2 [T]
x/tool2 [T] imports: x/tool [T]`,
		}, {
			name:        "standard-config-standard-proj",
			givenRoot:   "standard-proj",
			givenConfig: `{"tool": ["x/*"], "db": ["db/*"]}`,
			expectedDeps: `/ [G] imports: db/store [D]
/ [G] imports: domain1 [S]
/ [G] imports: domain2 [S]
db/store [D] imports: db/model [D]
db/store [D] imports: x/tool [T]
db/store [D] imports: x/tool2 [T]
domain1 [S] imports: db/store [D]
domain1 [S] imports: x/tool [T]
domain2 [S] imports: db/store [D]
domain2 [S] imports: x/tool2 [T]`,
		}, {
			name:        "standard-config-complex-proj",
			givenRoot:   "complex-proj",
			givenConfig: `{"tool": ["pkg/x/*"], "db": ["pkg/db/*"]}`,
			expectedDeps: `cmd/exe1 [G] imports: pkg/db/store [D]
cmd/exe1 [G] imports: pkg/domain1 [S]
cmd/exe1 [G] imports: pkg/domain2 [S]
cmd/exe2 [G] imports: pkg/db/store [D]
cmd/exe2 [G] imports: pkg/domain3 [S]
cmd/exe2 [G] imports: pkg/domain4 [S]
pkg/db/store [D] imports: pkg/db/model [D]
pkg/db/store [D] imports: pkg/x/tool [T]
pkg/db/store [D] imports: pkg/x/tool2 [T]
pkg/domain1 [S] imports: pkg/db/store [D]
pkg/domain1 [S] imports: pkg/x/tool [T]
pkg/domain2 [S] imports: pkg/db/store [D]
pkg/domain2 [S] imports: pkg/x/tool2 [T]
pkg/domain3 [S] imports: pkg/db/store [D]
pkg/domain3 [S] imports: pkg/x/tool [T]
pkg/domain4 [S] imports: pkg/db/store [D]
pkg/domain4 [S] imports: pkg/domain3 [S]
pkg/domain4 [S] imports: pkg/x/tool2 [T]`,
		}, {
			name:        "standard-config-half-pkgs-proj",
			givenRoot:   "half-pkgs-proj",
			givenConfig: `{tool: ["x/*"], db: ["db/*"]}`,
			expectedDeps: `/ [G] imports: db/store [D]
/ [G] imports: domain1 [S]
/ [G] imports: domain2 [S]
db/store [D] imports: db/model [D]
db/store [D] imports: db/store/substore [S]
db/store [D] imports: x/tool [T]
db/store [D] imports: x/tool2 [T]
domain1 [S] imports: db/store [D]
domain1 [S] imports: x/tool [T]
domain2 [S] imports: db/store [D]
domain2 [S] imports: x/tool2 [T]
x/tool [T] imports: x/tool/subtool [S]`,
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

func failWithDiff(t *testing.T, expected, actual string) {
	exps := strings.Split(expected, "\n")
	acts := strings.Split(actual, "\n")

	i := 0
	j := 0
	n := len(exps) - 1
	m := len(acts) - 1
	if n >= 0 && exps[n] == "" {
		n--
	}
	if m >= 0 && acts[m] == "" {
		m--
	}
	for i <= n && j <= m {
		if exps[i] < acts[j] {
			t.Errorf("expected but missing:  %s", exps[i])
			i++
		} else if exps[i] == acts[j] {
			i++
			j++
		} else if exps[i] > acts[j] {
			t.Errorf("actual but unexpected: %s", acts[j])
			j++
		}
	}
	for ; i <= n; i++ {
		t.Errorf("expected but missing:  %s", exps[i])
	}
	for ; j <= m; j++ {
		t.Errorf("actual but unexpected: %s", acts[j])
	}
}

func sortedImpNames(imps map[string]data.PkgType) []string {
	names := make([]string, 0, len(imps))
	for imp := range imps {
		names = append(names, imp)
	}
	sort.Strings(names)
	return names
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

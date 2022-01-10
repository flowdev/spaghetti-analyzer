package analdata_test

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-cutter/data"
)

func TestDependencyMapSortedPkgNames(t *testing.T) {
	givenDepMap := map[string]analdata.PkgImports{
		"b":  analdata.PkgImports{},
		"a":  analdata.PkgImports{},
		"aa": analdata.PkgImports{},
		"b ": analdata.PkgImports{},
		"a ": analdata.PkgImports{},
	}
	expectedNames := []string{"a", "a ", "aa", "b", "b "}

	actualNames := analdata.DependencyMap(givenDepMap).SortedPkgNames()
	if len(actualNames) != len(givenDepMap) {
		t.Errorf("expected %d names, actual %d",
			len(givenDepMap), len(actualNames))
	}
	if !reflect.DeepEqual(actualNames, expectedNames) {
		t.Errorf("expected names to be %q, got %q", expectedNames, actualNames)
	}
}

func TestFilterDepMap(t *testing.T) {
	manyLinks, err := data.NewSimplePatternList([]string{"a", "b/c/d", "e**", "f/**"}, "test")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	bigDepMap := analdata.DependencyMap{
		"a": analdata.PkgImports{
			PkgType: data.TypeGod,
			Imports: map[string]data.PkgType{
				"b/c/d":   data.TypeTool,
				"epsilon": data.TypeTool,
				"escher":  data.TypeTool,
				"z":       data.TypeDB,
				"f":       data.TypeStandard,
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
				"f/j":    data.TypeStandard,
				"escher": data.TypeTool,
				"x":      data.TypeDB,
			},
		},
		"f/h": analdata.PkgImports{
			PkgType: data.TypeStandard,
			Imports: map[string]data.PkgType{
				"x": data.TypeDB,
				"m": data.TypeStandard,
			},
		},
		"f/i": analdata.PkgImports{
			PkgType: data.TypeStandard,
			Imports: map[string]data.PkgType{
				"escher": data.TypeTool,
			},
		},
	}
	t.Logf("bigDepMap:\n%s", prettyPrint(bigDepMap))

	specs := []struct {
		name           string
		givenIdx       int
		givenLinks     data.PatternList
		expectedDepMap string
	}{
		{
			name:           "negative-index",
			givenIdx:       -123,
			givenLinks:     manyLinks,
			expectedDepMap: prettyPrint(bigDepMap),
		}, {
			name:           "no-links",
			givenIdx:       0,
			givenLinks:     data.PatternList{},
			expectedDepMap: prettyPrint(bigDepMap),
		}, {
			name:           "no-wildcard-leave-package",
			givenIdx:       1,
			givenLinks:     manyLinks,
			expectedDepMap: ``,
		}, {
			name:       "no-wildcard-tree",
			givenIdx:   0,
			givenLinks: manyLinks,
			expectedDepMap: `a [G] imports: b/c/d [T]
a [G] imports: epsilon [T]
a [G] imports: escher [T]
a [G] imports: f [S]
a [G] imports: z [D]
x [D] imports: b/c/d [T]
x [D] imports: escher [T]
z [D] imports: b/c/d [T]
z [D] imports: epsilon [T]
z [D] imports: escher [T]
z [D] imports: x [D]`,
		}, {
			name:           "wildcard-leave-packages",
			givenIdx:       2,
			givenLinks:     manyLinks,
			expectedDepMap: ``,
		}, {
			name:       "wildcard-tree",
			givenIdx:   3,
			givenLinks: manyLinks,
			expectedDepMap: `f [S] imports: f/g [S]
f [S] imports: f/h [S]
f [S] imports: f/i [S]
f/g [S] imports: escher [T]
f/g [S] imports: f/j [S]
f/g [S] imports: x [D]
f/h [S] imports: m [S]
f/h [S] imports: x [D]
f/i [S] imports: escher [T]
m [S] imports: b/c/d [T]
m [S] imports: x [D]
x [D] imports: b/c/d [T]
x [D] imports: escher [T]`,
		},
	}

	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualDepMap := analdata.FilterDepMap(bigDepMap, spec.givenIdx, spec.givenLinks)
			sDeps := prettyPrint(actualDepMap)
			if sDeps != spec.expectedDepMap {
				failWithDiff(t, spec.expectedDepMap, sDeps)
			}
		})
	}
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

func TestPkgForPattern(t *testing.T) {
	specs := []struct {
		name         string
		givenPattern string
		expectedPkg  string
	}{
		{
			name:         "no-wildcard",
			givenPattern: "abc/def/ghi",
			expectedPkg:  "abc/def/ghi",
		}, {
			name:         "one-simple-wildcard",
			givenPattern: "abc/def/*",
			expectedPkg:  "abc/def",
		}, {
			name:         "one-middle-wildcard",
			givenPattern: "abc/def/gh*i",
			expectedPkg:  "abc/def",
		}, {
			name:         "many-wildcards",
			givenPattern: "abc/d**e*f/*",
			expectedPkg:  "abc",
		}, {
			name:         "no-slash",
			givenPattern: "abc*",
			expectedPkg:  "",
		},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualPkg := analdata.PkgForPattern(spec.givenPattern)
			if actualPkg != spec.expectedPkg {
				t.Errorf("expected package %q, actual %q", spec.expectedPkg, actualPkg)
			}
		})
	}
}

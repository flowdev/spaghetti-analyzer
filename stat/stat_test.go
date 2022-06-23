package stat_test

import (
	"strings"
	"testing"

	"github.com/pkg/diff"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-analyzer/stat"
	"github.com/flowdev/spaghetti-cutter/data"
)

var bigDepMap = analdata.DependencyMap{
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

func TestGenerateTable(t *testing.T) {
	specs := []struct {
		name         string
		givenDepMap  analdata.DependencyMap
		expectedStat string
	}{
		{
			name:         "empty",
			givenDepMap:  analdata.DependencyMap{},
			expectedStat: "",
		}, {
			name: "minimal",
			givenDepMap: analdata.DependencyMap{
				"a": analdata.PkgImports{
					PkgType: data.TypeGod,
					Imports: map[string]data.PkgType{
						"b": data.TypeTool,
					},
				},
			},
			expectedStat: stat.Header +
				`| [a](#package-a) | [ \[G\] ](#legend) | [1](#direct-dependencies-imports-of-package-a) | [1](#all-including-transitive-dependencies-imports-of-package-a) | 0 | 0 | 0 |
` + stat.Legend + `

### Package a


#### Direct Dependencies (Imports) Of Package a
` + "`b`" + `

#### All (Including Transitive) Dependencies (Imports) Of Package a
` + "`b`" + `
`,
		}, {
			name:        "complex",
			givenDepMap: bigDepMap,
			expectedStat: stat.Header +
				`| [a](#package-a) | [ \[G\] ](#legend) | [5](#direct-dependencies-imports-of-package-a) | [11](#all-including-transitive-dependencies-imports-of-package-a) | 0 | 0 | 0 |
| [f](#package-f) | [ \[S\] ](#legend) | [3](#direct-dependencies-imports-of-package-f) | [8](#all-including-transitive-dependencies-imports-of-package-f) | [1](#packages-using-importing-package-f) | [5](#packages-shielded-from-users-of-package-f) | [5](#packages-shielded-from-all-users-of-package-f) |
| [f/g](#package-fg) | [ \[S\] ](#legend) | [3](#direct-dependencies-imports-of-package-fg) | [4](#all-including-transitive-dependencies-imports-of-package-fg) | [1](#packages-using-importing-package-fg) | [1](#packages-shielded-from-users-of-package-fg) | [1](#packages-shielded-from-all-users-of-package-fg) |
| [f/h](#package-fh) | [ \[S\] ](#legend) | [2](#direct-dependencies-imports-of-package-fh) | [4](#all-including-transitive-dependencies-imports-of-package-fh) | [1](#packages-using-importing-package-fh) | [1](#packages-shielded-from-users-of-package-fh) | [1](#packages-shielded-from-all-users-of-package-fh) |
| [f/i](#package-fi) | [ \[S\] ](#legend) | [1](#direct-dependencies-imports-of-package-fi) | [1](#all-including-transitive-dependencies-imports-of-package-fi) | [1](#packages-using-importing-package-fi) | 0 | 0 |
| [m](#package-m) | [ \[S\] ](#legend) | [2](#direct-dependencies-imports-of-package-m) | [3](#all-including-transitive-dependencies-imports-of-package-m) | [1](#packages-using-importing-package-m) | 0 | 0 |
| [x](#package-x) | [ \[D\] ](#legend) | [2](#direct-dependencies-imports-of-package-x) | [2](#all-including-transitive-dependencies-imports-of-package-x) | [4](#packages-using-importing-package-x) | [3](#packages-shielded-from-users-of-package-x) | 0 |
| [z](#package-z) | [ \[D\] ](#legend) | [4](#direct-dependencies-imports-of-package-z) | [4](#all-including-transitive-dependencies-imports-of-package-z) | [1](#packages-using-importing-package-z) | 0 | 0 |
` + stat.Legend + `

### Package a


#### Direct Dependencies (Imports) Of Package a
` + "`b/c/d`, `epsilon`, `escher`" + `, [f](#package-f), [z](#package-z)

#### All (Including Transitive) Dependencies (Imports) Of Package a
` + "`b/c/d`, `epsilon`, `escher`, [f](#package-f), [f/g](#package-fg), [f/h](#package-fh), [f/i](#package-fi), `f/j`" + `, [m](#package-m), [x](#package-x), [z](#package-z)

### Package f


#### Direct Dependencies (Imports) Of Package f
[f/g](#package-fg), [f/h](#package-fh), [f/i](#package-fi)

#### All (Including Transitive) Dependencies (Imports) Of Package f
` + "`b/c/d`, `escher`, [f/g](#package-fg), [f/h](#package-fh), [f/i](#package-fi), `f/j`" + `, [m](#package-m), [x](#package-x)

#### Packages Using (Importing) Package f
[a](#package-a)

#### Packages Shielded From Users Of Package f
* [a](#package-a): [f/g](#package-fg), [f/h](#package-fh), [f/i](#package-fi), ` + "`f/j`" + `, [m](#package-m)


#### Packages Shielded From All Users Of Package f
[f/g](#package-fg), [f/h](#package-fh), [f/i](#package-fi), ` + "`f/j`" + `, [m](#package-m)

### Package f/g


#### Direct Dependencies (Imports) Of Package f/g
` + "`escher`, `f/j`" + `, [x](#package-x)

#### All (Including Transitive) Dependencies (Imports) Of Package f/g
` + "`b/c/d`, `escher`, `f/j`" + `, [x](#package-x)

#### Packages Using (Importing) Package f/g
[f](#package-f)

#### Packages Shielded From Users Of Package f/g
* [f](#package-f): ` + "`f/j`" + `


#### Packages Shielded From All Users Of Package f/g
` + "`f/j`" + `

### Package f/h


#### Direct Dependencies (Imports) Of Package f/h
[m](#package-m), [x](#package-x)

#### All (Including Transitive) Dependencies (Imports) Of Package f/h
` + "`b/c/d`, `escher`" + `, [m](#package-m), [x](#package-x)

#### Packages Using (Importing) Package f/h
[f](#package-f)

#### Packages Shielded From Users Of Package f/h
* [f](#package-f): [m](#package-m)


#### Packages Shielded From All Users Of Package f/h
[m](#package-m)

### Package f/i


#### Direct Dependencies (Imports) Of Package f/i
` + "`escher`" + `

#### All (Including Transitive) Dependencies (Imports) Of Package f/i
` + "`escher`" + `

#### Packages Using (Importing) Package f/i
[f](#package-f)

### Package m


#### Direct Dependencies (Imports) Of Package m
` + "`b/c/d`" + `, [x](#package-x)

#### All (Including Transitive) Dependencies (Imports) Of Package m
` + "`b/c/d`, `escher`" + `, [x](#package-x)

#### Packages Using (Importing) Package m
[f/h](#package-fh)

### Package x


#### Direct Dependencies (Imports) Of Package x
` + "`b/c/d`, `escher`" + `

#### All (Including Transitive) Dependencies (Imports) Of Package x
` + "`b/c/d`, `escher`" + `

#### Packages Using (Importing) Package x
[f/g](#package-fg), [f/h](#package-fh), [m](#package-m), [z](#package-z)

#### Packages Shielded From Users Of Package x
* [f/g](#package-fg): ` + "`b/c/d`" + `
* [f/h](#package-fh): ` + "`escher`" + `
* [m](#package-m): ` + "`escher`" + `


### Package z


#### Direct Dependencies (Imports) Of Package z
` + "`b/c/d`, `epsilon`, `escher`" + `, [x](#package-x)

#### All (Including Transitive) Dependencies (Imports) Of Package z
` + "`b/c/d`, `epsilon`, `escher`" + `, [x](#package-x)

#### Packages Using (Importing) Package z
[a](#package-a)
`,
		},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualStat := stat.Generate(spec.givenDepMap)
			if actualStat != spec.expectedStat {
				t.Error(diffStats(t, spec.expectedStat, actualStat))
			}
		})
	}
}

func diffStats(t *testing.T, want, got string) string {
	buf := &strings.Builder{}

	err := diff.Text("expected", "got", want, got, buf)
	if err != nil {
		t.Fatalf("unable to diff result: %v", err)
	}
	return buf.String()
}

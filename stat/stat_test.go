package stat_test

import (
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-analyzer/stat"
	"github.com/flowdev/spaghetti-cutter/data"
	"github.com/pkg/diff"
)

// "github.com/pkg/diff"

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

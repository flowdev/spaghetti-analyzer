package doc_test

import (
	"strings"
	"testing"

	"github.com/pkg/diff"

	"github.com/flowdev/spaghetti-analyzer/analdata"
	"github.com/flowdev/spaghetti-analyzer/doc"
	"github.com/flowdev/spaghetti-cutter/data"
)

func TestGenerateTable(t *testing.T) {
	specs := []struct {
		name          string
		givenIdx      int
		givenLinks    data.PatternList
		givenDocFiles []string
		givenDepMap   analdata.DependencyMap
		expectedDoc   string
	}{
		{
			name:          "minimal",
			givenIdx:      0,
			givenLinks:    newPatternList(t, "a"),
			givenDocFiles: []string{"./" + doc.FileName},
			givenDepMap: analdata.DependencyMap{
				"a": analdata.PkgImports{
					PkgType: data.TypeGod,
					Imports: map[string]data.PkgType{
						"b": data.TypeTool,
					},
				},
			},
			expectedDoc: doc.Title + `github.com/flowdev/tst/a

| | b - T |
| :- | :- |
| **a** | **T** |
` + doc.Legend,
		},
	}
	rootPkg := "github.com/flowdev/tst"
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualDoc := doc.GenerateTable(spec.givenIdx, spec.givenLinks, spec.givenDocFiles, spec.givenDepMap, rootPkg)
			if actualDoc != spec.expectedDoc {
				t.Error(diffDocs(t, spec.expectedDoc, actualDoc))
			}
		})
	}
}

func newPatternList(t *testing.T, pkgs ...string) data.PatternList {
	pl, err := data.NewSimplePatternList(pkgs, "test")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	return pl
}

func diffDocs(t *testing.T, want, got string) string {
	buf := &strings.Builder{}

	err := diff.Text("expected", "got", want, got, buf)
	if err != nil {
		t.Fatalf("unable to diff result: %v", err)
	}
	return buf.String()
}

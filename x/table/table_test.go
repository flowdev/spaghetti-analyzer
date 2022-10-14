package table_test

import (
	"strings"
	"testing"

	"github.com/flowdev/spaghetti-analyzer/x/table"
	"github.com/pkg/diff"
)

func TestGenerate(t *testing.T) {
	specs := []struct {
		name          string
		givenData     [][]string
		givenAlign    []table.Align
		expectedTable string
	}{
		{
			name:          "nil",
			givenData:     nil,
			givenAlign:    nil,
			expectedTable: "",
		}, {
			name:          "empty",
			givenData:     [][]string{},
			givenAlign:    []table.Align{},
			expectedTable: "",
		}, {
			name: "minimal",
			givenData: [][]string{
				{"heading"},
				{"cell value"},
			},
			givenAlign: []table.Align{table.AlignCenter},
			expectedTable: `|  heading   |
|:----------:|
| cell value |
`,
		}, {
			name: "complex",
			givenData: [][]string{
				{"heading 1", "", "heading 3", "heading 4"},
				{"cell value 1", "", "value 3"},
			},
			givenAlign: []table.Align{
				table.AlignLeft, table.AlignCenter, table.AlignRight, table.AlignLeft, table.AlignLeft,
			},
			expectedTable: `| heading 1    |   | heading 3 | heading 4 |
|:-------------|:-:|----------:|:----------|
| cell value 1 |   |   value 3 |           |
`,
		},
	}
	for _, spec := range specs {
		t.Run(spec.name, func(t *testing.T) {
			actualTable := table.Generate(spec.givenData, spec.givenAlign)
			if actualTable != spec.expectedTable {
				t.Error(diffTables(t, spec.expectedTable, actualTable))
			}
		})
	}
}

func diffTables(t *testing.T, want, got string) string {
	buf := &strings.Builder{}

	err := diff.Text("expected", "got", want, got, buf)
	if err != nil {
		t.Fatalf("unable to diff result: %v", err)
	}
	return buf.String()
}

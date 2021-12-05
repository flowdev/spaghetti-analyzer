package tool2_test

import (
	"testing"

	"github.com/flowdev/spaghetti-analyzer/deps/testdata/half-pkgs-proj/x/tool2"
)

func TestTool(t *testing.T) {
	t.Log("Executing TestTool")
	tool2.Tool2()
}

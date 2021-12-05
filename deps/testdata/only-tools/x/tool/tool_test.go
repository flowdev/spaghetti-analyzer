package tool_test

import (
	"testing"

	"github.com/flowdev/spaghetti-analyzer/deps/testdata/only-tools/x/tool"
)

func TestTool(t *testing.T) {
	t.Log("Executing TestTool")
	tool.Tool()
}

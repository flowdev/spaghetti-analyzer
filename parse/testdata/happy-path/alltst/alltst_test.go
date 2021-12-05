package alltst_test

import (
	"testing"

	"github.com/flowdev/spaghetti-analyzer/parse/testdata/happy-path/alltst"
)

func TestAlltst(t *testing.T) {
	t.Log("Executing TestAlltst")
	alltst.Alltst()
}

package loop

import (
	"testing"

	"github.com/vasiliyyudin/go-mutesting/test"
)

func TestMutatorLoopCondition(t *testing.T) {
	test.Mutator(
		t,
		MutatorLoopCondition,
		"../../testdata/loop/condition.go",
		2,
	)
}

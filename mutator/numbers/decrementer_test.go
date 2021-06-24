package numbers

import (
	"testing"

	"github.com/vasiliyyudin/go-mutesting/test"
)

func TestMutatorNumbersDecrementer(t *testing.T) {
	test.Mutator(
		t,
		MutatorNumbersDecrementer,
		"../../testdata/numbers/decrementer.go",
		2,
	)
}

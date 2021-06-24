package arithmetic

import (
	"testing"

	"github.com/vasiliyyudin/go-mutesting/test"
)

func TestMutatorArithmeticAssignment(t *testing.T) {
	test.Mutator(
		t,
		MutatorArithmeticAssignment,
		"../../testdata/arithmetic/assignment.go",
		11,
	)
}

package conditional

import (
	"testing"

	"github.com/vasiliyyudin/go-mutesting/test"
)

func TestMutatorConditionalNegated(t *testing.T) {
	test.Mutator(
		t,
		MutatorConditionalNegated,
		"../../testdata/conditional/negated.go",
		6,
	)
}

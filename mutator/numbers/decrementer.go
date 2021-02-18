package numbers

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"github.com/avito-tech/go-mutesting/mutator"
)

func init() {
	mutator.Register("numbers/decrementer", MutatorNumbersDecrementer)
}

// MutatorNumbersDecrementer implements a mutator to decrement int and float.
func MutatorNumbersDecrementer(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.BasicLit)
	if !ok {
		return nil
	}

	if n.Kind == token.INT {
		original := n.Value
		originalInt, err := strconv.Atoi(n.Value)
		if err != nil {
			return nil
		}

		originalInt--
		mutated := strconv.Itoa(originalInt)

		return []mutator.Mutation{
			{
				Change: func() {
					n.Value = mutated
				},
				Reset: func() {
					n.Value = original
				},
			},
		}
	}

	if n.Kind == token.FLOAT {
		original := n.Value
		originalFloat, err := strconv.ParseFloat(n.Value, 64)
		if err != nil {
			return nil
		}

		originalFloat--
		mutated := strconv.FormatFloat(originalFloat, 'f', -1, 64)

		return []mutator.Mutation{
			{
				Change: func() {
					n.Value = mutated
				},
				Reset: func() {
					n.Value = original
				},
			},
		}
	}

	return nil
}

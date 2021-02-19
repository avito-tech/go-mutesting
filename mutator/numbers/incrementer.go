package numbers

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"github.com/avito-tech/go-mutesting/mutator"
)

func init() {
	mutator.Register("numbers/incrementer", MutatorNumbersIncrementer)
}

// MutatorNumbersIncrementer implements a mutator to increment int and float.
func MutatorNumbersIncrementer(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
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

		originalInt++
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

		originalFloat++
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

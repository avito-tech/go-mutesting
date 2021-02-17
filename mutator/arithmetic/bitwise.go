package arithmetic

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/avito-tech/go-mutesting/mutator"
)

func init() {
	mutator.Register("arithmetic/bitwise", MutatorArithmeticBitwise)
}

var bitwiseMutations = map[token.Token]token.Token{
	token.AND:     token.OR,
	token.OR:      token.AND,
	token.XOR:     token.AND,
	token.AND_NOT: token.AND,
	token.SHL:     token.SHR,
	token.SHR:     token.SHL,
}

// MutatorArithmeticBitwise implements a mutator to change bitwise arithmetic.
func MutatorArithmeticBitwise(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	original := n.Op
	mutated, ok := bitwiseMutations[n.Op]
	if !ok {
		return nil
	}

	return []mutator.Mutation{
		{
			Change: func() {
				n.Op = mutated
			},
			Reset: func() {
				n.Op = original
			},
		},
	}
}

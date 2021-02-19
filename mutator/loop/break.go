package loop

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/avito-tech/go-mutesting/mutator"
)

func init() {
	mutator.Register("loop/break", MutatorLoopBreak)
}

var breakMutations = map[token.Token]token.Token{
	token.CONTINUE: token.BREAK,
	token.BREAK:    token.CONTINUE,
}

// MutatorLoopBreak implements a mutator to change continue to break and break to continue.
func MutatorLoopBreak(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.BranchStmt)
	if !ok {
		return nil
	}

	original := n.Tok
	mutated, ok := breakMutations[n.Tok]
	if !ok {
		return nil
	}

	return []mutator.Mutation{
		{
			Change: func() {
				n.Tok = mutated
			},
			Reset: func() {
				n.Tok = original
			},
		},
	}
}

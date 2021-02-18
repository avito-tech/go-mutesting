package loop

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/avito-tech/go-mutesting/mutator"
)

func init() {
	mutator.Register("loop/condition", MutatorLoopCondition)
}

// MutatorLoopCondition implements a mutator to change loop condition to always false.
func MutatorLoopCondition(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.ForStmt)
	if !ok {
		return nil
	}

	condition, ok := n.Cond.(*ast.BinaryExpr)
	if !ok {
		return nil
	}

	originalX := condition.X
	originalOp := condition.Op
	originalY := condition.Y

	return []mutator.Mutation{
		{
			Change: func() {
				condition.X = ast.NewIdent("1")
				condition.Op = token.LSS
				condition.Y = ast.NewIdent("1")
			},
			Reset: func() {
				condition.X = originalX
				condition.Op = originalOp
				condition.Y = originalY
			},
		},
	}
}

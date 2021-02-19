package loop

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/avito-tech/go-mutesting/mutator"
)

func init() {
	mutator.Register("loop/range_break", MutatorLoopRangeBreak)
}

// MutatorLoopRangeBreak implements a mutator to add a break to range-loop body.
func MutatorLoopRangeBreak(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.RangeStmt)
	if !ok {
		return nil
	}

	newBody := &ast.BlockStmt{
		List: []ast.Stmt{},
	}
	oldBody := n.Body

	newBreakStmt := &ast.BranchStmt{Tok: token.BREAK}
	newBody.List = append(newBody.List, newBreakStmt)

	for _, nn := range n.Body.List {
		newBody.List = append(newBody.List, nn)
	}

	return []mutator.Mutation{
		{
			Change: func() {
				n.Body = newBody
			},
			Reset: func() {
				n.Body = oldBody
			},
		},
	}
}

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

func checkIsRangeStatement(node ast.Stmt) bool {
	switch node.(type) {
	case *ast.RangeStmt:
		return true
	}

	return false
}

// MutatorLoopRangeBreak implements a mutator to add a break to range-loop body.
func MutatorLoopRangeBreak(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
	var list []ast.Stmt

	switch n := node.(type) {
	case *ast.BlockStmt:
		list = n.List
	case *ast.CaseClause:
		list = n.Body
	default:
		return nil
	}

	var mutations []mutator.Mutation

	for i, ni := range list {
		if checkIsRangeStatement(ni) {
			listIndex := i
			oldNode := list[listIndex]
			rangeNode, ok := oldNode.(*ast.RangeStmt)
			if !ok {
				return nil
			}

			newBody := &ast.BlockStmt{
				List: []ast.Stmt{},
			}

			newBreakStmt := &ast.BranchStmt{Tok: token.BREAK}
			newBody.List = append(newBody.List, newBreakStmt)

			for _, nn := range rangeNode.Body.List {
				newBody.List = append(newBody.List, nn)
			}

			newNode := &ast.RangeStmt{
				For:    rangeNode.For,
				Key:    rangeNode.Key,
				Value:  rangeNode.Value,
				TokPos: rangeNode.TokPos,
				Tok:    rangeNode.Tok,
				Body:   newBody,
				X:      rangeNode.X,
			}

			mutations = append(mutations, mutator.Mutation{
				Change: func() {
					list[listIndex] = newNode
				},
				Reset: func() {
					list[listIndex] = oldNode
				},
			})
		}
	}

	return mutations
}

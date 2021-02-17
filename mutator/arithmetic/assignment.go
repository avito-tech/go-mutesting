package arithmetic

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/avito-tech/go-mutesting/mutator"
)

func init() {
	mutator.Register("arithmetic/assignment", MutatorArithmeticAssignment)
}

var assignmentMutations = map[token.Token]token.Token{
	token.ADD_ASSIGN:     token.ASSIGN,
	token.SUB_ASSIGN:     token.ASSIGN,
	token.MUL_ASSIGN:     token.ASSIGN,
	token.QUO_ASSIGN:     token.ASSIGN,
	token.REM_ASSIGN:     token.ASSIGN,
	token.AND_ASSIGN:     token.ASSIGN,
	token.OR_ASSIGN:      token.ASSIGN,
	token.XOR_ASSIGN:     token.ASSIGN,
	token.SHL_ASSIGN:     token.ASSIGN,
	token.SHR_ASSIGN:     token.ASSIGN,
	token.AND_NOT_ASSIGN: token.ASSIGN,
}

// MutatorArithmeticAssignment implements a mutator to change base assign logic.
func MutatorArithmeticAssignment(_ *types.Package, _ *types.Info, node ast.Node) []mutator.Mutation {
	n, ok := node.(*ast.AssignStmt)
	if !ok {
		return nil
	}

	original := n.Tok
	mutated, ok := assignmentMutations[n.Tok]
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

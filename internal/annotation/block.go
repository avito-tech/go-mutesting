package annotation

import (
	"go/ast"
	"go/token"
)

var StatNodesInBlockForRegex = make(map[int]map[token.Pos]MutatorInfo)
var StatNodesInBlockForLine = make(map[int]map[token.Pos]MutatorInfo)

// HandleBlockStmt is a temporary workaround specifically for handling BlockStmt nodes in AST.
// It performs cleanup and transfers collected annotation data to statement nodes within blocks.
// This is a tactical solution to handle edge cases where mutators only look at nodes inside block statements.
// A more robust architectural solution should be implemented in future versions.
func HandleBlockStmt(node ast.Stmt) bool {
	for _, n := range StatNodesInBlockForRegex {
		if mutatorName, exists := n[node.Pos()]; exists {
			if shouldSkipMutator(mutatorName, "statement/remove") {
				return true
			}
		}
	}

	for _, n := range StatNodesInBlockForLine {
		if mutatorName, exists := n[node.Pos()]; exists {
			if shouldSkipMutator(mutatorName, "statement/remove") {
				return true
			}
		}
	}

	return false
}

func CleanupGlobalStatBlock() {
	StatNodesInBlockForRegex = make(map[int]map[token.Pos]MutatorInfo)
	StatNodesInBlockForLine = make(map[int]map[token.Pos]MutatorInfo)
}

func (r *RegexAnnotation) CopyToStatNodesInBlock() {
	for line, nodes := range r.Exclusions {
		if _, exists := StatNodesInBlockForRegex[line]; !exists {
			StatNodesInBlockForRegex[line] = make(map[token.Pos]MutatorInfo)
		}

		for pos, mutatorInfo := range nodes {
			StatNodesInBlockForRegex[line][pos] = mutatorInfo
		}
	}
}

func (l *LineAnnotation) CopyToStatNodesInBlock() {
	for line, nodes := range l.Exclusions {
		if _, exists := StatNodesInBlockForLine[line]; !exists {
			StatNodesInBlockForLine[line] = make(map[token.Pos]MutatorInfo)
		}

		for pos, mutatorInfo := range nodes {
			StatNodesInBlockForLine[line][pos] = mutatorInfo
		}
	}
}

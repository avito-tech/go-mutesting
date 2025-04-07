package annotation

import (
	"go/ast"
	"go/token"
)

var statNodesInBlockForRegex = make(map[int]map[token.Pos]mutatorInfo)
var statNodesInBlockForLine = make(map[int]map[token.Pos]mutatorInfo)

// HandleBlockStmt is a temporary workaround specifically for handling BlockStmt nodes in AST.
// It performs cleanup and transfers collected annotation data to statement nodes within blocks.
// This is a tactical solution to handle edge cases where mutators only look at nodes inside block statements.
// A more robust architectural solution should be implemented in future versions.
func HandleBlockStmt(node ast.Stmt) bool {
	for _, n := range statNodesInBlockForRegex {
		if mutatorName, exists := n[node.Pos()]; exists {
			if shouldSkipMutator(mutatorName, "statement/remove") {
				return true
			}
		}
	}

	for _, n := range statNodesInBlockForLine {
		if mutatorName, exists := n[node.Pos()]; exists {
			if shouldSkipMutator(mutatorName, "statement/remove") {
				return true
			}
		}
	}

	return false
}

func cleanupGlobalStatBlock() {
	statNodesInBlockForRegex = make(map[int]map[token.Pos]mutatorInfo)
	statNodesInBlockForLine = make(map[int]map[token.Pos]mutatorInfo)
}

func (r *RegexAnnotation) copyToStatNodesInBlock() {
	for line, nodes := range r.Exclusions {
		if _, exists := statNodesInBlockForRegex[line]; !exists {
			statNodesInBlockForRegex[line] = make(map[token.Pos]mutatorInfo)
		}

		for pos, mutatorInfo := range nodes {
			statNodesInBlockForRegex[line][pos] = mutatorInfo
		}
	}
}

func (l *LineAnnotation) copyToStatNodesInBlock() {
	for line, nodes := range l.Exclusions {
		if _, exists := statNodesInBlockForLine[line]; !exists {
			statNodesInBlockForLine[line] = make(map[token.Pos]mutatorInfo)
		}

		for pos, mutatorInfo := range nodes {
			statNodesInBlockForLine[line][pos] = mutatorInfo
		}
	}
}

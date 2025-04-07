package annotation

import (
	"go/ast"
	"go/token"
	"strings"
)

type LineAnnotation struct {
	Exclusions map[int]map[token.Pos]MutatorInfo
	Name       string
}

// parseLineAnnotation parses a comment line containing a next-line annotation.
func (l *LineAnnotation) parseLineAnnotation(comment string) MutatorInfo {
	content := strings.TrimSpace(strings.TrimPrefix(comment, l.Name))
	if content == "" {
		return MutatorInfo{}
	}

	mutators := parseMutators(content)

	return MutatorInfo{
		Names: mutators,
	}
}

// collectNodesOnNextLine processes a "mutator-disable-next-line" annotation.
// It:
// 1. Parses the mutator names from the annotation comment
// 2. Determines the line number immediately following the comment
// 3. Collects all AST nodes that appear on that line
// 4. Records the exclusion information for those nodes
func (l *LineAnnotation) collectNodesOnNextLine(comment *ast.Comment, fset *token.FileSet, file *ast.File) {
	mutators := l.parseLineAnnotation(comment.Text)

	start, end := findLine(fset, comment)
	var nextLine int
	if start == end {
		nextLine = start + 1
	}

	lines := []int{nextLine}

	collectExcludedNodes(fset, file, lines, l.Exclusions, mutators)
}

// filterNodesOnNextLine checks if a given node should be excluded from mutation based on:
// 1. Whether the node appears in the Exclusions map
// 2. Whether the current mutator is in the node's exclusion list
func (l *LineAnnotation) filterNodesOnNextLine(node ast.Node, mutatorName string) bool {
	for _, nnn := range l.Exclusions {
		if mutatorInfo, exists := nnn[node.Pos()]; exists {
			if shouldSkipMutator(mutatorInfo, mutatorName) {
				return true
			}
		}
	}

	return false
}

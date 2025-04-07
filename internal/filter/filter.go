package filter

import (
	"go/ast"
	"go/token"
)

type NodeCollector interface {
	Collect(file *ast.File, fset *token.FileSet, fileAbs string)
}

type NodeFilter interface {
	ShouldSkip(node ast.Node, mutatorName string) bool
}

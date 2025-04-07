package annotation

import (
	"go/ast"
	"go/token"
	"strings"
)

type FunctionAnnotation struct {
	Exclusions map[token.Pos]struct{}
	Name       string
}

// collectFunctions records all nodes within a function declaration to be excluded from mutation.
// It collects both the function declaration itself and all its child nodes.
func (f *FunctionAnnotation) collectFunctions(fun *ast.FuncDecl) {
	f.Exclusions[fun.Pos()] = struct{}{}

	ast.Inspect(fun, func(n ast.Node) bool {
		if n != nil {
			f.Exclusions[n.Pos()] = struct{}{}
		}

		return true
	})
}

// filterFunctions checks whether a given node should be excluded from mutation
func (f *FunctionAnnotation) filterFunctions(node ast.Node) bool {
	if _, exists := f.Exclusions[node.Pos()]; exists {
		return true
	}

	return false
}

// existsFuncAnnotation checks if a function declaration has the annotation
func (p *Processor) existsFuncAnnotation(f *ast.FuncDecl) bool {
	if f.Doc == nil {
		return false
	}

	for _, comment := range f.Doc.List {
		if strings.HasPrefix(comment.Text, p.FunctionAnnotation.Name) {
			return true
		}
	}

	return false
}

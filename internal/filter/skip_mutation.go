package filter

import (
	"go/ast"
	"go/token"
)

// SkipMakeFilter is a filter that tracks numeric arguments in 'make' calls
// for slices and maps to be ignored during mutation.
type SkipMakeFilter struct {
	// IgnoredNodes maps positions of numeric literals to their parent 'make' call expressions
	IgnoredNodes map[token.Pos]*ast.CallExpr
}

func NewMakeSkipper() *SkipMakeFilter {
	return &SkipMakeFilter{IgnoredNodes: make(map[token.Pos]*ast.CallExpr)}
}

// Collect collects numeric arguments (children) from 'make' calls (parents) for slices/maps to be ignored during mutation
func (s *SkipMakeFilter) Collect(file *ast.File, _ *token.FileSet, _ string) {
	ast.Inspect(file, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "make" && len(callExpr.Args) > 1 {
				arg0 := callExpr.Args[0]
				_, isArray := arg0.(*ast.ArrayType)
				_, isMap := arg0.(*ast.MapType)
				if isArray || isMap {
					if lit, ok := callExpr.Args[1].(*ast.BasicLit); ok && lit.Kind == token.INT {
						s.IgnoredNodes[lit.Pos()] = callExpr
					}
					if len(callExpr.Args) > 2 {
						if lit, ok := callExpr.Args[2].(*ast.BasicLit); ok && lit.Kind == token.INT {
							s.IgnoredNodes[lit.Pos()] = callExpr
						}
					}
					return false
				}
			}
		}
		return true
	})
}

// ShouldSkip determines whether a given AST node should be skipped during mutation.
func (s *SkipMakeFilter) ShouldSkip(node ast.Node, _ string) bool {
	_, exists := s.IgnoredNodes[node.Pos()]
	return exists
}

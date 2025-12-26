package filter

import (
	"go/ast"
	"go/token"
)

// SkipMakeArgsFilter is a filter that tracks numeric arguments in 'make' calls
// for slices and maps to be ignored during mutation.
type SkipMakeArgsFilter struct {
	// IgnoredNodes maps positions of numeric literals to their parent 'make' call expressions
	IgnoredNodes map[token.Pos]*ast.CallExpr
}

// NewSkipMakeArgsFilter creates and returns a new initialized Processor.
func NewSkipMakeArgsFilter() *SkipMakeArgsFilter {
	return &SkipMakeArgsFilter{IgnoredNodes: make(map[token.Pos]*ast.CallExpr)}
}

// Collect collects numeric arguments (children) from 'make' calls (parents) for slices/maps to be ignored during mutation
func (s *SkipMakeArgsFilter) Collect(file *ast.File, _ *token.FileSet, _ string) {
	ast.Inspect(file, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "make" && len(callExpr.Args) > 1 {
				arg0 := callExpr.Args[0]
				_, isArray := arg0.(*ast.ArrayType)
				_, isMap := arg0.(*ast.MapType)
				_, isIdent := arg0.(*ast.Ident)
				if isArray || isMap || isIdent {
					s.collectForIgnoredNodes(callExpr.Args[1], callExpr)

					if len(callExpr.Args) > 2 {
						s.collectForIgnoredNodes(callExpr.Args[2], callExpr)
					}
					return false
				}
			}
		}
		return true
	})
}

// ShouldSkip determines whether a given AST node should be skipped during mutation.
func (s *SkipMakeArgsFilter) ShouldSkip(node ast.Node, _ string) bool {
	_, exists := s.IgnoredNodes[node.Pos()]
	return exists
}

// collectForIgnoredNodes recursively collects all numeric literals and unary/binary operators in an expression
func (s *SkipMakeArgsFilter) collectForIgnoredNodes(expr ast.Expr, callExpr *ast.CallExpr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// Direct numeric literal
		if e.Kind == token.INT {
			s.IgnoredNodes[e.Pos()] = callExpr
		}
	case *ast.BinaryExpr:
		// Binary operations (addition, subtraction, multiplication, division, etc.)
		s.IgnoredNodes[e.OpPos] = callExpr
		s.collectForIgnoredNodes(e.X, callExpr)
		s.collectForIgnoredNodes(e.Y, callExpr)
	case *ast.CallExpr:
		// Calling a function (e.g. len())
		for _, arg := range e.Args {
			s.collectForIgnoredNodes(arg, callExpr)
		}
	case *ast.ParenExpr:
		// Expression in brackets
		s.collectForIgnoredNodes(e.X, callExpr)
	case *ast.UnaryExpr:
		// Unary operators (+, -, etc.)
		if xLit, ok := e.X.(*ast.BasicLit); ok && xLit.Kind == token.INT {
			// If a unary operator is applied to a numeric literal, both the operator and the literal itself are added.
			s.IgnoredNodes[e.OpPos] = callExpr
			s.IgnoredNodes[xLit.Pos()] = callExpr
		} else {
			// Otherwise, we continue the tour inside.
			s.collectForIgnoredNodes(e.X, callExpr)
		}
	}
}

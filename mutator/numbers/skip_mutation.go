package numbers

import (
	"go/ast"
	"go/token"
)

// to associate positions of numeric literals with their corresponding 'make' call expressions
var ignoredNodes = make(map[token.Pos]*ast.CallExpr)

// collects numeric arguments (children) from 'make' calls (parents) for slices/maps to be ignored during mutation
func skipMutationForMakeArgs(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "make" && len(callExpr.Args) > 1 {
				arg0 := callExpr.Args[0]
				_, isArray := arg0.(*ast.ArrayType)
				_, isMap := arg0.(*ast.MapType)
				if isArray || isMap {
					if lit, ok := callExpr.Args[1].(*ast.BasicLit); ok && lit.Kind == token.INT {
						ignoredNodes[lit.Pos()] = callExpr
					}
					if len(callExpr.Args) > 2 {
						if lit, ok := callExpr.Args[2].(*ast.BasicLit); ok && lit.Kind == token.INT {
							ignoredNodes[lit.Pos()] = callExpr
						}
					}
					return false
				}
			}
		}
		return true
	})
}

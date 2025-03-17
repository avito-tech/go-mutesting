package numbers

import (
	"go/ast"
	"go/token"
)

var globalParentMap = make(map[token.Pos]*ast.CallExpr)

func skipMutationForMake(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "make" && len(callExpr.Args) > 1 {
				arg0 := callExpr.Args[0]
				_, isArray := arg0.(*ast.ArrayType)
				_, isMap := arg0.(*ast.MapType)
				if isArray || isMap {
					if lit, ok := callExpr.Args[1].(*ast.BasicLit); ok && lit.Kind == token.INT {
						globalParentMap[lit.Pos()] = callExpr
					}
					if len(callExpr.Args) > 2 {
						if lit, ok := callExpr.Args[2].(*ast.BasicLit); ok && lit.Kind == token.INT {
							globalParentMap[lit.Pos()] = callExpr
						}
					}
					return false
				}
			}
		}
		return true
	})
}

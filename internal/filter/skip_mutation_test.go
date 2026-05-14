package filter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipMutationForInitSlicesAndMaps(t *testing.T) {
	tests := []struct {
		name              string
		code              string
		expectedLiterals  []string
		expectedOperators []string
	}{
		{
			name:              "skip mutation for slice init with len",
			code:              `package main; var a = make([]int, 10)`,
			expectedLiterals:  []string{"10"},
			expectedOperators: []string{},
		},
		{
			name:              "skip mutation for slice init with len and cap",
			code:              `package main; var a = make([]int, 10, 20)`,
			expectedLiterals:  []string{"10", "20"},
			expectedOperators: []string{},
		},
		{
			name: "skip mutation for slice init in assigning inside a struct",
			code: `package ccc; 
				   type TestCase struct { 
					Devices    []DeviceStatus 
				   }; 
				   type DeviceStatus struct {
					DeviceName     string
					Status         string 
					ReportViewerID string 
				   }; 
				   func fff() { 
					testCase := &TestCase{ Devices: make([]DeviceStatus, 0) } 
				   }`,
			expectedLiterals:  []string{"0"},
			expectedOperators: []string{},
		},
		{
			name:              "skip mutation for map init with cap",
			code:              `package main; var a = make(map[int]bool, 0)`,
			expectedLiterals:  []string{"0"},
			expectedOperators: []string{},
		},
		{
			name:              "do not skip mutation for slice init with variable",
			code:              `package main; var x = 10; var a = make([]int, x)`,
			expectedLiterals:  []string{},
			expectedOperators: []string{},
		},
		{
			name:              "do not skip mutation for other literals",
			code:              `package main; var a = 42`,
			expectedLiterals:  []string{},
			expectedOperators: []string{},
		},
		{
			name:              "skip mutation for slice with type alias",
			code:              `package main; type Contents []*string; var a = make(Contents, 5)`,
			expectedLiterals:  []string{"5"},
			expectedOperators: []string{},
		},
		{
			name:              "skip mutation for complex expression with binary op",
			code:              `package main; var a = make([]int, 0, len(arr)+1)`,
			expectedLiterals:  []string{"0", "1"},
			expectedOperators: []string{"+"},
		},
		{
			name:              "skip mutation for unary operations",
			code:              `package main; var a = make([]int, -5, +10)`,
			expectedLiterals:  []string{"5", "10"},
			expectedOperators: []string{"-", "+"},
		},
		{
			name:              "skip mutation for complex nested expressions",
			code:              `package main; var a = make([]int, (3), len(arr)+2*4)`,
			expectedLiterals:  []string{"3", "2", "4"},
			expectedOperators: []string{"+", "*"},
		},
		{
			// kills mutants that relax len(callExpr.Args) > 1 to >= 1, > 0, or true:
			// with those mutants the 1-arg make(MySlice) enters the isIdent branch
			// and panics on Args[1] access.
			name:              "do not skip mutation for single-arg make with type alias",
			code:              `package main; type MySlice []int; var a = make(MySlice)`,
			expectedLiterals:  []string{},
			expectedOperators: []string{},
		},
		{
			// kills mutants that remove or noop the else branch of the UnaryExpr case:
			// without recursion into -(2+3), the "2", "3", and "+" inside are not collected.
			name:              "skip mutation for unary negation of complex inner expression",
			code:              `package main; var a = make([]int, 0, -(2+3))`,
			expectedLiterals:  []string{"0", "2", "3"},
			expectedOperators: []string{"+"},
		},
		{
			// kills mutants that break or remove the collectForIgnoredNodes call inside the
			// CallExpr case loop: with those mutants only the first arg (or nothing) is collected.
			name:              "skip mutation for nested call with multiple int args",
			code:              `package main; var a = make([]int, someFunc(2, 3))`,
			expectedLiterals:  []string{"2", "3"},
			expectedOperators: []string{},
		},
		{
			// kills the mutant that replaces xLit.Kind == token.INT with true: that mutant
			// treats float unary operands like int ones and adds their operator positions.
			name:              "do not skip mutation for unary on float literals",
			code:              `package main; var a = make([]float64, -1.5, +2.3)`,
			expectedLiterals:  []string{},
			expectedOperators: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, err := parser.ParseFile(fs, "skip_mutation_test.go", tt.code, parser.Mode(0))
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			s := NewSkipMakeArgsFilter()
			s.Collect(node, fs, "")

			var foundLiterals []string
			var foundOperators []string

			ast.Inspect(node, func(n ast.Node) bool {
				if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.INT {
					if _, exists := s.IgnoredNodes[lit.Pos()]; exists {
						foundLiterals = append(foundLiterals, lit.Value)
					}
				}

				if binExpr, ok := n.(*ast.BinaryExpr); ok {
					if _, exists := s.IgnoredNodes[binExpr.OpPos]; exists {
						foundOperators = append(foundOperators, binExpr.Op.String())
					}
				}

				if unaryExpr, ok := n.(*ast.UnaryExpr); ok {
					if _, exists := s.IgnoredNodes[unaryExpr.OpPos]; exists {
						foundOperators = append(foundOperators, unaryExpr.Op.String())
					}
				}

				return true
			})

			assert.ElementsMatch(t, tt.expectedLiterals, foundLiterals, "Literals mismatch")
			assert.ElementsMatch(t, tt.expectedOperators, foundOperators, "Operators mismatch")
		})
	}
}

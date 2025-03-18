package numbers

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestSkipMutationForInitSlicesAndMaps(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "skip mutation for slice init with len",
			code:     `package main; var a = make([]int, 10)`,
			expected: true,
		},
		{
			name:     "skip mutation for slice init with len and cap",
			code:     `package main; var a = make([]int, 10, 20)`,
			expected: true,
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
			expected: true,
		},
		{
			name:     "skip mutation for map init with cap",
			code:     `package main; var a = make(map[int]bool, 0)`,
			expected: true,
		},
		{
			name:     "do not skip mutation for slice init with variable",
			code:     `package main; var x = 10; var a = make([]int, x)`,
			expected: false,
		},
		{
			name:     "do not skip mutation for other literals",
			code:     `package main; var a = 42`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, err := parser.ParseFile(fs, "skip_mutation_test.go", tt.code, parser.Mode(0))
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			ignoredNodes = make(map[token.Pos]*ast.CallExpr)
			skipMutationForMakeArgs(node)

			var result bool
			ast.Inspect(node, func(n ast.Node) bool {
				if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.INT {
					_, found := ignoredNodes[lit.Pos()]
					result = found
					return false
				}
				return true
			})

			assert.Equal(t, tt.expected, result)
		})
	}
}

package filter

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

func TestCoverageFilter(t *testing.T) {
	t.Parallel()

	covData := `mode: set
github.com/example/pkg/test.go:1.1,2.2 1 1
github.com/example/pkg/test.go:3.3,4.4 1 0
`
	tmpDir := t.TempDir()
	covFile := filepath.Join(tmpDir, "coverage.out")

	err := os.WriteFile(covFile, []byte(covData), 0o644)
	if err != nil {
		t.Fatalf("failed to write cov file: %v", err)
	}

	cf, err := NewCoverageFilter(covFile)
	if err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}

	if cf == nil {
		t.Fatal("expected non-nil filter")
	}

	// Line 1 and 2 should be true
	if !cf.CoveredLines["github.com/example/pkg/test.go"][1] {
		t.Error("expected line 1 to be covered")
	}

	if !cf.CoveredLines["github.com/example/pkg/test.go"][2] {
		t.Error("expected line 2 to be covered")
	}

	// Line 3 and 4 should be absent or false
	if cf.CoveredLines["github.com/example/pkg/test.go"][3] {
		t.Error("expected line 3 to be not covered")
	}

	fset := token.NewFileSet()
	src := `package main
func main() { // Line 2
	println("covered") // Line 3
	println("not covered") // Line 4
}
`

	f, err := parser.ParseFile(fset, "github.com/example/pkg/test.go", src, 0)
	if err != nil {
		t.Fatalf("failed to parse file: %v", err)
	}

	cf.Collect(f, fset, "")

	funcDecl, ok := f.Decls[0].(*ast.FuncDecl)
	if !ok {
		t.Fatalf("expected ast.FuncDecl, got %T", f.Decls[0])
	}

	body := funcDecl.Body

	// The nodes will be at line 3 and 4 in the parsed AST.
	// Since our dummy coverage says line 3 has 0 coverage, ShouldSkip should return true.
	if !cf.ShouldSkip(body.List[0], "") {
		t.Error("expected line 3 to be skipped")
	}

	if !cf.ShouldSkip(body.List[1], "") {
		t.Error("expected line 4 to be skipped")
	}
}

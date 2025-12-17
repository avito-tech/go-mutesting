package annotation

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/avito-tech/go-mutesting/internal/filter"
	"github.com/avito-tech/go-mutesting/mutator"
)

// Annotation constants define the comment patterns used to disable mutations
const (
	FuncAnnotation     = "// mutator-disable-func"
	RegexpAnnotation   = "// mutator-disable-regexp"
	NextLineAnnotation = "// mutator-disable-next-line"
)

// Processor handles mutation exclusion logic based on source code annotations.
type Processor struct {
	options

	FunctionAnnotation FunctionAnnotation
	RegexAnnotation    RegexAnnotation
	LineAnnotation     LineAnnotation
}

// NewProcessor creates and returns a new initialized Processor.
func NewProcessor(optionFunc ...OptionFunc) *Processor {
	opts := options{}

	for _, f := range optionFunc {
		f(&opts)
	}

	processor := &Processor{
		options: opts,
		FunctionAnnotation: FunctionAnnotation{
			Exclusions: make(map[token.Pos]struct{}), // *ast.FuncDecl node + all its children
			Name:       FuncAnnotation},
		RegexAnnotation: RegexAnnotation{
			GlobalRegexCollector: NewRegexCollector(opts.global.filteredRegexps),
			Exclusions:           make(map[int]map[token.Pos]mutatorInfo), // source code line -> node -> excluded mutators
			Name:                 RegexpAnnotation,
		},
		LineAnnotation: LineAnnotation{
			Exclusions: make(map[int]map[token.Pos]mutatorInfo), // source code line -> node -> excluded mutators
			Name:       NextLineAnnotation,
		},
	}

	return processor
}

type mutatorInfo struct {
	Names []string
}

// Collect processes an AST file to gather all mutation exclusions based on annotations.
func (p *Processor) Collect(
	file *ast.File,
	fset *token.FileSet,
	fileAbs string,
) {
	// comment based collectors
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if p.existsFuncAnnotation(f) {
				p.FunctionAnnotation.collectFunctions(f)
			}
		}
	}

	handler := p.buildChain()

	for _, commentGroup := range file.Comments {
		for _, comm := range commentGroup.List {
			name := getAnnotationName(comm)
			handler.Handle(name, comm, fset, file, fileAbs)
		}
	}

	p.RegexAnnotation.GlobalRegexCollector.Collect(fset, file, fileAbs)

	p.collectNodesForBlockStmt()
}

// ShouldSkip determines if a given node should be excluded from mutation.
func (p *Processor) ShouldSkip(node ast.Node, mutatorName string) bool {
	return p.FunctionAnnotation.filterFunctions(node) ||
		p.RegexAnnotation.filterRegexNodes(node, mutatorName) ||
		p.LineAnnotation.filterNodesOnNextLine(node, mutatorName)
}

// DecoratorFilter creates a mutator that applies one or more filters before executing the provided mutator.
func DecoratorFilter(m mutator.Mutator, name string, filters ...filter.NodeFilter) mutator.Mutator {
	return func(pkg *types.Package, info *types.Info, node ast.Node) []mutator.Mutation {
		for _, f := range filters {
			if f.ShouldSkip(node, name) {
				return nil
			}
		}

		return m(pkg, info, node)
	}
}

// getAnnotationName identifies the type of annotation
func getAnnotationName(comment *ast.Comment) string {
	content := strings.TrimSpace(comment.Text)
	if strings.HasPrefix(content, RegexpAnnotation) {
		return RegexpAnnotation
	}
	if strings.HasPrefix(content, NextLineAnnotation) {
		return NextLineAnnotation
	}
	if strings.HasPrefix(content, FuncAnnotation) {
		return FuncAnnotation
	}

	return ""
}

// collectExcludedNodes populates a map of nodes to exclude from mutation based on the line numbers.
func collectExcludedNodes(
	fileSet *token.FileSet,
	file *ast.File,
	lines []int,
	excludedNodes map[int]map[token.Pos]mutatorInfo,
	mutators mutatorInfo) {
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		startLine, endLine := getNodeLineRange(fileSet, n)

		for _, line := range lines {
			if startLine == line || endLine == line {
				if _, exists := excludedNodes[line]; !exists {
					excludedNodes[line] = make(map[token.Pos]mutatorInfo)
				}
				excludedNodes[line][n.Pos()] = mutators
			}
		}

		return true
	})
}

// collectNodesForBlockStmt is a temporary workaround specifically for handling BlockStmt nodes in AST.
// It performs cleanup and transfers collected annotation data to statement nodes within blocks.
// This is a tactical solution to handle edge cases where mutators only look at nodes inside block statements.
// A more robust architectural solution should be implemented in future versions.
func (p *Processor) collectNodesForBlockStmt() {
	cleanupGlobalStatBlock()
	p.RegexAnnotation.copyToStatNodesInBlock()
	p.LineAnnotation.copyToStatNodesInBlock()
}

// parseMutators parses a comma-separated string of mutator names into a clean slice of strings.
func parseMutators(mutatorList string) []string {
	mutators := make([]string, 0)

	rawTargets := strings.Split(mutatorList, ",")
	for _, t := range rawTargets {
		name := strings.TrimSpace(t)
		if name != "" {
			mutators = append(mutators, name)
		}
	}

	return mutators
}

// shouldSkipMutator determines whether a specific mutator should be skipped
func shouldSkipMutator(mutatorInfo mutatorInfo, mutatorName string) bool {
	for _, name := range mutatorInfo.Names {
		if name == mutatorName || name == "*" {
			return true
		}
	}

	return false
}

// getNodeLineRange calculates the line number range (start to end) that a given AST node occupies in the source file.
func getNodeLineRange(fileSet *token.FileSet, node ast.Node) (startLine, endLine int) {
	startPos := fileSet.Position(node.Pos())
	endPos := fileSet.Position(node.End())

	return startPos.Line, endPos.Line
}

// findLine determines the line number range of a comment node.
func findLine(fileSet *token.FileSet, comment *ast.Comment) (int, int) {
	startLine, endLine := getNodeLineRange(fileSet, comment)

	return startLine, endLine
}

package annotation

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSkipMutator(t *testing.T) {
	tests := []struct {
		name        string
		mutatorInfo mutatorInfo
		mutatorName string
		expected    bool
	}{
		{
			name:        "Mutator matches exactly",
			mutatorInfo: mutatorInfo{Names: []string{"MutatorA", "MutatorB"}},
			mutatorName: "MutatorA",
			expected:    true,
		},
		{
			name:        "Mutator matches wildcard",
			mutatorInfo: mutatorInfo{Names: []string{"*"}},
			mutatorName: "AnyMutator",
			expected:    true,
		},
		{
			name:        "Mutator does not match",
			mutatorInfo: mutatorInfo{Names: []string{"MutatorA", "MutatorB"}},
			mutatorName: "MutatorC",
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkipMutator(tt.mutatorInfo, tt.mutatorName)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestParseMutators(t *testing.T) {
	tests := []struct {
		name        string
		mutatorList string
		expected    []string
	}{
		{
			name:        "Valid list of mutators",
			mutatorList: "MutatorA, MutatorB, MutatorC",
			expected:    []string{"MutatorA", "MutatorB", "MutatorC"},
		},
		{
			name:        "List with leading and trailing spaces",
			mutatorList: "  MutatorA,  MutatorB , MutatorC  ",
			expected:    []string{"MutatorA", "MutatorB", "MutatorC"},
		},
		{
			name:        "Empty string",
			mutatorList: "",
			expected:    []string{},
		},
		{
			name:        "String with only commas",
			mutatorList: ",,,",
			expected:    []string{},
		},
		{
			name:        "Single mutator",
			mutatorList: "MutatorA",
			expected:    []string{"MutatorA"},
		},
		{
			name:        "Multiple empty elements",
			mutatorList: "MutatorA,,,MutatorB,,MutatorC,,",
			expected:    []string{"MutatorA", "MutatorB", "MutatorC"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMutators(tt.mutatorList)
			if !assert.Equal(t, result, tt.expected) {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestParseRegexAnnotation(t *testing.T) {
	tests := []struct {
		name          string
		commentText   string
		expectedRegex *regexp.Regexp
		expectedInfo  mutatorInfo
	}{
		{
			name:        "Valid regex and mutators",
			commentText: "RegexName ^[a-z]+$ MutatorA, MutatorB",
			expectedRegex: func() *regexp.Regexp {
				re, _ := regexp.Compile("^[a-z]+$")
				return re
			}(),
			expectedInfo: mutatorInfo{
				Names: []string{"MutatorA", "MutatorB"},
			},
		},
		{
			name:        "Valid regex without mutators",
			commentText: "RegexName ^[0-9]{4}$",
			expectedRegex: func() *regexp.Regexp {
				re, _ := regexp.Compile("^[0-9]{4}$")
				return re
			}(),
			expectedInfo: mutatorInfo{
				Names: []string{},
			},
		},
		{
			name:          "Invalid regex",
			commentText:   "RegexName [a-z",
			expectedRegex: nil,
			expectedInfo:  mutatorInfo{},
		},
		{
			name:          "Empty comment",
			commentText:   "RegexName ",
			expectedRegex: nil,
			expectedInfo:  mutatorInfo{},
		},
	}

	r := &RegexAnnotation{Name: "RegexName"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultRegex, resultInfo := r.parseRegexAnnotation(tt.commentText)

			if resultRegex == nil && tt.expectedRegex != nil || resultRegex != nil && resultRegex.String() != tt.expectedRegex.String() {
				t.Errorf("Expected regex %v, but got %v", tt.expectedRegex, resultRegex)
			}

			if len(resultInfo.Names) != len(tt.expectedInfo.Names) {
				t.Errorf("Expected mutators %v, but got %v", tt.expectedInfo.Names, resultInfo.Names)
			} else {
				for i, mutator := range resultInfo.Names {
					if mutator != tt.expectedInfo.Names[i] {
						t.Errorf("Expected mutator %v, but got %v", tt.expectedInfo.Names[i], mutator)
					}
				}
			}
		})
	}
}

func TestParseLineAnnotation(t *testing.T) {
	tests := []struct {
		name         string
		commentText  string
		expectedInfo mutatorInfo
	}{
		{
			name:        "Valid mutators",
			commentText: "LineName MutatorA, MutatorB, MutatorC",
			expectedInfo: mutatorInfo{
				Names: []string{"MutatorA", "MutatorB", "MutatorC"},
			},
		},
		{
			name:        "Single mutator",
			commentText: "LineName MutatorA",
			expectedInfo: mutatorInfo{
				Names: []string{"MutatorA"},
			},
		},
		{
			name:        "Multiple mutators with spaces",
			commentText: "LineName  MutatorA, MutatorB , MutatorC  ",
			expectedInfo: mutatorInfo{
				Names: []string{"MutatorA", "MutatorB", "MutatorC"},
			},
		},
		{
			name:        "Empty comment",
			commentText: "LineName ",
			expectedInfo: mutatorInfo{
				Names: []string{},
			},
		},
		{
			name:        "Only spaces in mutators",
			commentText: "LineName ,,,",
			expectedInfo: mutatorInfo{
				Names: []string{},
			},
		},
		{
			name:        "Empty mutators",
			commentText: "LineName",
			expectedInfo: mutatorInfo{
				Names: []string{},
			},
		},
	}

	l := &LineAnnotation{Name: "LineName"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := l.parseLineAnnotation(tt.commentText)

			if len(result.Names) != len(tt.expectedInfo.Names) {
				t.Errorf("Expected mutators %v, but got %v", tt.expectedInfo.Names, result.Names)
			} else {
				for i, mutator := range result.Names {
					if mutator != tt.expectedInfo.Names[i] {
						t.Errorf("Expected mutator %v, but got %v", tt.expectedInfo.Names[i], mutator)
					}
				}
			}
		})
	}
}

func TestExistsFuncAnnotation(t *testing.T) {
	tests := []struct {
		name          string
		funcDecl      *ast.FuncDecl
		expectedExist bool
	}{
		{
			name: "No annotations",
			funcDecl: &ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "// Some other comment"},
					},
				},
			},
			expectedExist: false,
		},
		{
			name: "Valid annotation exists",
			funcDecl: &ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "// mutator-disable-func something here"},
					},
				},
			},
			expectedExist: true,
		},
		{
			name: "Multiple comments, valid annotation exists",
			funcDecl: &ast.FuncDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "// Another comment"},
						{Text: "// mutator-disable-func something here"},
					},
				},
			},
			expectedExist: true,
		},
		{
			name: "Doc is nil",
			funcDecl: &ast.FuncDecl{
				Doc: nil,
			},
			expectedExist: false,
		},
	}

	p := NewProcessor()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.existsFuncAnnotation(tt.funcDecl)
			if result != tt.expectedExist {
				t.Errorf("Expected %v, but got %v", tt.expectedExist, result)
			}
		})
	}
}

func TestCollectFunctionsAndFilterFunctions(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		expected  bool
		filterPos token.Pos
	}{
		{
			name:      "Function with one statement",
			code:      `package main; func test() { var a = 10 }`,
			expected:  true,
			filterPos: token.Pos(1),
		},
		{
			name:      "Function with nested statements",
			code:      `package main; func test() { if true { var a = 10 } }`,
			expected:  true,
			filterPos: token.Pos(2),
		},
	}

	f := &FunctionAnnotation{Exclusions: make(map[token.Pos]struct{})}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := token.NewFileSet()
			node, err := parser.ParseFile(fs, "func_annotation_test.go", tt.code, parser.Mode(0))
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			// Находим первую функцию в коде
			var funcDecl *ast.FuncDecl
			ast.Inspect(node, func(n ast.Node) bool {
				if fDecl, ok := n.(*ast.FuncDecl); ok {
					funcDecl = fDecl
					return false
				}
				return true
			})

			f.collectFunctions(funcDecl)

			filtered := f.filterFunctions(funcDecl)
			assert.Equal(t, tt.expected, filtered)

		})
	}
}

func TestFindLinesMatchedRegex(t *testing.T) {
	tests := []struct {
		name          string
		filePath      string
		re            *regexp.Regexp
		expectedLines []int
	}{
		{
			name:          "No regex match",
			filePath:      "../../testdata/annotation/regex.go",
			re:            regexp.MustCompile("notmatching"),
			expectedLines: []int{},
		},
		{
			name:          "Match variable declaration",
			filePath:      "../../testdata/annotation/regex.go",
			re:            regexp.MustCompile(`test`),
			expectedLines: []int{37, 38},
		},
		{
			name:          "Match Println",
			filePath:      "../../testdata/annotation/regex.go",
			re:            regexp.MustCompile(`Println`),
			expectedLines: []int{16, 17, 23, 33, 38},
		},
		{
			name:          "Multiple matches on multiple lines",
			filePath:      "../../testdata/annotation/regex.go",
			re:            regexp.MustCompile(`xx+`),
			expectedLines: []int{12, 13, 16, 17},
		},
		{
			name:          "Match MyStruct",
			filePath:      "../../testdata/annotation/regex.go",
			re:            regexp.MustCompile(`MyStruct`),
			expectedLines: []int{19, 27, 31},
		},
		{
			name:          "Match interface declaration",
			filePath:      "../../testdata/annotation/regex.go",
			re:            regexp.MustCompile(`interface`),
			expectedLines: []int{42},
		},
		{
			name:          "Match slog.Info",
			filePath:      "../../testdata/annotation/regex.go",
			re:            regexp.MustCompile(`slog\.Info`),
			expectedLines: []int{49},
		},
		{
			name:          "Match s.Method()",
			filePath:      "../../testdata/annotation/regex.go",
			re:            regexp.MustCompile(`s\.Method\(\)`),
			expectedLines: []int{20},
		},
		{
			name:          "Regex is nil",
			filePath:      "../../testdata/annotation/regex.go",
			re:            nil,
			expectedLines: []int{},
		},
		{
			name:          "Empty file",
			filePath:      "../../testdata/annotation/empty.go",
			re:            regexp.MustCompile(`test`),
			expectedLines: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := &RegexAnnotation{Name: "TestRegexAnnotation"}

			actual, _ := r.findLinesMatchingRegex(tt.filePath, tt.re)

			assert.ElementsMatch(t, tt.expectedLines, actual)
		})
	}
}

func TestCollect(t *testing.T) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "../../testdata/annotation/collect.go", nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		t.Fatalf("failed to parse file: %v", err)
	}

	processor := NewProcessor()

	processor.Collect(file, fs, "../../testdata/annotation/collect.go")

	assert.NotEmpty(t, processor.FunctionAnnotation.Exclusions)
	assert.Equal(t, processor.FunctionAnnotation.Exclusions, map[token.Pos]struct{}{
		75: {}, 99: {}, 104: {}, 114: {}, 115: {}, 117: {}, 122: {}, 126: {}, 129: {}, 136: {}, 140: {},
	})

	assert.NotEmpty(t, processor.RegexAnnotation.Exclusions)
	assert.Equal(t, processor.RegexAnnotation.Exclusions, map[int]map[token.Pos]mutatorInfo{
		14: {
			169: {Names: []string{"*"}},
			173: {Names: []string{"*"}},
			181: {Names: []string{"*"}},
		},
		22: {
			304: {Names: []string{"*"}},
			308: {Names: []string{"*"}},
			316: {Names: []string{"*"}},
		},
		21: {
			288: {Names: []string{"*"}},
			292: {Names: []string{"*"}},
			300: {Names: []string{"*"}},
		},
	})

	assert.NotEmpty(t, processor.LineAnnotation.Exclusions)
	assert.Equal(t, processor.LineAnnotation.Exclusions, map[int]map[token.Pos]mutatorInfo{
		19: {
			275: {Names: []string{"numbers/incrementer"}},
			279: {Names: []string{"numbers/incrementer"}},
			283: {Names: []string{"numbers/incrementer"}},
		},
	})

}

func TestCollectGlobal(t *testing.T) {
	filePath := "../../testdata/annotation/global/collect.go"

	fs := token.NewFileSet()
	file, err := parser.ParseFile(
		fs,
		filePath,
		nil,
		parser.AllErrors|parser.ParseComments,
	)
	assert.NoError(t, err)

	processor := NewProcessor(WithGlobalRegexpFilter("\\.Log *"))

	processor.Collect(file, fs, filePath)

	assert.NotEmpty(t, processor.RegexAnnotation.GlobalRegexCollector.Exclusions)
	assert.Equal(t, processor.RegexAnnotation.GlobalRegexCollector.Exclusions, map[int]map[token.Pos]mutatorInfo{
		17: {
			256: {Names: []string{"*"}},
			263: {Names: []string{"*"}},
			267: {Names: []string{"*"}},
		},
	})
}

package annotation

import (
	"bufio"
	"go/ast"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
)

type RegexAnnotation struct {
	Exclusions map[int]map[token.Pos]mutatorInfo
	Name       string
}

// parseRegexAnnotation parses a comment line containing a regex annotation.
func (r *RegexAnnotation) parseRegexAnnotation(comment string) (*regexp.Regexp, mutatorInfo) {
	content := strings.TrimSpace(strings.TrimPrefix(comment, r.Name))
	if content == "" {
		return nil, mutatorInfo{}
	}

	parts := strings.SplitN(content, " ", 2)

	pattern := strings.TrimSpace(parts[0])
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Warning: invalid regex in annotation: %q, error: %v\n", pattern, err)
		return nil, mutatorInfo{}
	}

	mutators := make([]string, 0)
	if len(parts) > 1 {
		mutators = parseMutators(parts[1])
	}

	return re, mutatorInfo{
		Names: mutators,
	}
}

// collectMatchNodes processes a "mutator-disable-regexp" annotation comment by:
// 1. Parsing the regex pattern and mutators from the comment
// 2. Finding all lines in the file that match the regex
// 3. Recording nodes from matching lines to be excluded
func (r *RegexAnnotation) collectMatchNodes(comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string) {
	regex, mutators := r.parseRegexAnnotation(comment.Text)

	lines, err := r.findLinesMatchingRegex(fileAbs, regex)
	if err != nil {
		log.Printf("Error scaning a source file: %v", err)
	}

	collectExcludedNodes(fset, file, lines, r.Exclusions, mutators)
}

// findLinesMatchingRegex scans a source file and returns line numbers that match the given regex.
func (r *RegexAnnotation) findLinesMatchingRegex(filePath string, regex *regexp.Regexp) ([]int, error) {
	var matchedLineNumbers []int

	if regex == nil {
		return matchedLineNumbers, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}

	reader := bufio.NewReader(f)

	lineNumber := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if regex.MatchString(line) {
			matchedLineNumbers = append(matchedLineNumbers, lineNumber+1)
		}
		lineNumber++
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Error while file closing duting processing regex annotation: %v", err.Error())
		}
	}()

	return matchedLineNumbers, nil
}

// filterRegexNodes checks if a given node should be excluded from mutation based on:
// 1. Whether the node appears in the Exclusions map
// 2. Whether the current mutator is in the node's exclusion list
func (r *RegexAnnotation) filterRegexNodes(node ast.Node, mutatorName string) bool {
	for _, nodes := range r.Exclusions {
		if mutatorInfo, exists := nodes[node.Pos()]; exists {
			if shouldSkipMutator(mutatorInfo, mutatorName) {
				return true
			}
		}
	}

	return false
}

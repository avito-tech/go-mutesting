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

// Collector defines the interface for handlers.
// Implementations should handle specific annotation types.
type Collector interface {
	// Handle processes an annotation if it matches the handler's type,
	// otherwise delegates to the next handler in the chain.
	Handle(name string, comment *ast.Comment, fset *token.FileSet, file *ast.File, fileAbs string)
}

// RegexExclusion structure that contains info required for ast.Node exclusion from mutations
type RegexExclusion struct {
	regex    *regexp.Regexp
	mutators mutatorInfo
}

// RegexCollector Collector based on regular expressions parse all file
type RegexCollector struct {
	Exclusions            map[int]map[token.Pos]mutatorInfo
	GlobalExclusionsRegex []RegexExclusion
}

// NewRegexCollector constructor for RegexCollector
func NewRegexCollector(
	exclusionsConfig []string,
) RegexCollector {
	exclusionsRegex := make([]RegexExclusion, 0, len(exclusionsConfig))
	for _, exclusion := range exclusionsConfig {
		re, inf := parseConfig(exclusion)
		if re != nil {
			exclusionsRegex = append(exclusionsRegex, RegexExclusion{
				regex:    re,
				mutators: inf,
			})
		}
	}

	return RegexCollector{
		Exclusions:            make(map[int]map[token.Pos]mutatorInfo),
		GlobalExclusionsRegex: exclusionsRegex,
	}
}

// Collect processes regex pattern
func (r *RegexCollector) Collect(
	fset *token.FileSet,
	file *ast.File,
	fileAbs string,
) {
	for _, exclusions := range r.GlobalExclusionsRegex {
		lines, err := r.findLinesMatchingRegex(fileAbs, exclusions.regex)
		if err != nil {
			log.Printf("Error scaning a source file: %v", err)
		}

		if len(lines) > 0 {
			collectExcludedNodes(fset, file, lines, r.Exclusions, exclusions.mutators)
		}
	}
}

func parseConfig(configLine string) (*regexp.Regexp, mutatorInfo) {
	// splitted[0] - contains regexp splitted[1] contains mutators
	splitted := strings.SplitN(configLine, " ", 2)

	if len(splitted) < 1 {
		return nil, mutatorInfo{}
	}

	pattern := splitted[0]
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Warning: invalid regex in annotation: %q, error: %v\n", pattern, err)
		return nil, mutatorInfo{}
	}

	var mutators []string
	if len(splitted) > 1 {
		mutators = parseMutators(splitted[1])
	} else {
		mutators = []string{"*"}
	}

	return re, mutatorInfo{
		Names: mutators,
	}
}

// findLinesMatchingRegex scans a source file and returns line numbers that match the given regex.
func (r *RegexCollector) findLinesMatchingRegex(filePath string, regex *regexp.Regexp) ([]int, error) {
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

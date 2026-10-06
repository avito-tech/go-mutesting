package filter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/cover"
)

// CoverageFilter filters out AST nodes that are on lines without test coverage.
type CoverageFilter struct {
	CoveredLines map[string]map[int]bool
	currentFset  *token.FileSet
}

// NewCoverageFilter creates a new CoverageFilter from a go cover profile.
func NewCoverageFilter(profilePath string) (*CoverageFilter, error) {
	cf := &CoverageFilter{
		CoveredLines: make(map[string]map[int]bool),
	}

	if profilePath == "" {
		return cf, nil
	}

	profiles, err := cover.ParseProfiles(profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse profiles: %w", err)
	}

	for _, p := range profiles {
		lines := make(map[int]bool)

		for _, b := range p.Blocks {
			if b.Count == 0 {
				continue
			}

			for i := b.StartLine; i <= b.EndLine; i++ {
				lines[i] = true
			}
		}

		cf.CoveredLines[p.FileName] = lines
	}

	return cf, nil
}

// Collect saves the fileset for the current file being walked.
func (c *CoverageFilter) Collect(_ *ast.File, fset *token.FileSet, _ string) {
	c.currentFset = fset
}

// ShouldSkip returns true if the node's line is not covered by tests.
func (c *CoverageFilter) ShouldSkip(node ast.Node, _ string) bool {
	if c.currentFset == nil || len(c.CoveredLines) == 0 {
		return false // No coverage filtering
	}

	pos := c.currentFset.Position(node.Pos())

	var bestMatch map[int]bool

	bestMatchLen := 0

	for k, v := range c.CoveredLines {
		aParts := strings.Split(k, "/")
		bParts := strings.Split(pos.Filename, "/")
		matchCount := 0

		for i := 1; i <= len(aParts) && i <= len(bParts); i++ {
			if aParts[len(aParts)-i] != bParts[len(bParts)-i] {
				break
			}

			matchCount++
		}

		if matchCount > bestMatchLen {
			bestMatchLen = matchCount
			bestMatch = v
		}
	}

	if bestMatchLen == 0 {
		return true // File not found in coverage data, skip it
	}

	return !bestMatch[pos.Line]
}

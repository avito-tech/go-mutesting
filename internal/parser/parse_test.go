package parser

import (
	"github.com/avito-tech/go-mutesting/internal/processor/annotation"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/avito-tech/go-mutesting/internal/filter"
)

func TestParseAndTypeCheckFileTypeCheckWholePackage(t *testing.T) {
	annotationProcessor := annotation.NewProcessor()
	skipFilterProcessor := filter.NewSkipMakeArgsFilter()

	collectors := []filter.NodeCollector{
		annotationProcessor,
		skipFilterProcessor,
	}
	_, _, _, _, err := ParseAndTypeCheckFile("../../astutil/create.go", collectors)
	assert.Nil(t, err)
}

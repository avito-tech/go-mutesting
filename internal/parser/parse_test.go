package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/avito-tech/go-mutesting/internal/annotation"
	"github.com/avito-tech/go-mutesting/internal/filter"
)

func TestParseAndTypeCheckFileTypeCheckWholePackage(t *testing.T) {
	annotationProcessor := annotation.NewProcessor()
	skipFilterProcessor := filter.NewMakeSkipper()

	collectors := []filter.NodeCollector{
		annotationProcessor,
		skipFilterProcessor,
	}
	_, _, _, _, err := ParseAndTypeCheckFile("../astutil/create.go", collectors)
	assert.Nil(t, err)
}

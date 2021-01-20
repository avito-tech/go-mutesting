package importing

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesOfArgs(t *testing.T) {
	p := os.Getenv("GOPATH") + "/src/"

	for _, test := range []struct {
		args   []string
		expect []string
	}{
		// empty
		{
			[]string{},
			[]string{"filepath.go", "import.go"},
		},
		// files
		{
			[]string{"filepath.go"},
			[]string{"filepath.go"},
		},
		// directories
		{
			[]string{"."},
			[]string{"filepath.go", "import.go"},
		},
		{
			[]string{".."},
			[]string{"../tool.go"},
		},
		{
			[]string{"../importing"},
			[]string{"../importing/filepath.go", "../importing/import.go"},
		},
		{
			[]string{"../..."},
			[]string{"../tool.go", "../importing/filepath.go", "../importing/import.go"},
		},
		// packages
		{
			[]string{"github.com/zimmski/go-tool"},
			[]string{p + "github.com/zimmski/go-tool/tool.go"},
		},
		{
			[]string{"github.com/zimmski/go-tool/importing"},
			[]string{p + "github.com/zimmski/go-tool/importing/filepath.go", p + "github.com/zimmski/go-tool/importing/import.go"},
		},
		{
			[]string{"github.com/zimmski/go-tool/..."},
			[]string{p + "github.com/zimmski/go-tool/tool.go", p + "github.com/zimmski/go-tool/importing/filepath.go", p + "github.com/zimmski/go-tool/importing/import.go"},
		},
	} {
		got := FilesOfArgs(test.args)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

func TestPackagesWithFilesOfArgs(t *testing.T) {
	p := os.Getenv("GOPATH") + "/src/"

	for _, test := range []struct {
		args   []string
		expect []Package
	}{
		// empty
		{
			[]string{},
			[]Package{Package{Name: ".", Files: []string{"filepath.go", "import.go"}}},
		},
		// files
		{
			[]string{"filepath.go"},
			[]Package{Package{Name: ".", Files: []string{"filepath.go"}}},
		},
		// directories
		{
			[]string{"."},
			[]Package{Package{Name: ".", Files: []string{"filepath.go", "import.go"}}},
		},
		{
			[]string{".."},
			[]Package{Package{Name: "..", Files: []string{"../tool.go"}}},
		},
		{
			[]string{"../importing"},
			[]Package{Package{Name: "../importing", Files: []string{"../importing/filepath.go", "../importing/import.go"}}},
		},
		{
			[]string{"../..."},
			[]Package{
				Package{Name: "..", Files: []string{"../tool.go"}},
				Package{Name: "../importing", Files: []string{"../importing/filepath.go", "../importing/import.go"}},
			},
		},
		// packages
		{
			[]string{"github.com/zimmski/go-tool"},
			[]Package{Package{Name: p + "github.com/zimmski/go-tool", Files: []string{p + "github.com/zimmski/go-tool/tool.go"}}},
		},
		{
			[]string{"github.com/zimmski/go-tool/importing"},
			[]Package{Package{Name: p + "github.com/zimmski/go-tool/importing", Files: []string{p + "github.com/zimmski/go-tool/importing/filepath.go", p + "github.com/zimmski/go-tool/importing/import.go"}}},
		},
		{
			[]string{"github.com/zimmski/go-tool/..."},
			[]Package{
				Package{Name: p + "github.com/zimmski/go-tool", Files: []string{p + "github.com/zimmski/go-tool/tool.go"}},
				Package{Name: p + "github.com/zimmski/go-tool/importing", Files: []string{p + "github.com/zimmski/go-tool/importing/filepath.go", p + "github.com/zimmski/go-tool/importing/import.go"}},
			},
		},
	} {
		got := PackagesWithFilesOfArgs(test.args)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

package importing

import (
	"fmt"
	"github.com/avito-tech/go-mutesting/internal/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesOfArgs(t *testing.T) {
	//TODO we need normal test with test folder
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
			[]string{"../importing"},
			[]string{"../importing/filepath.go", "../importing/import.go"},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal"},
			[]string(nil),
		},
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing"},
			[]string{p + "github.com/avito-tech/go-mutesting/internal/importing/filepath.go", p + "github.com/avito-tech/go-mutesting/internal/importing/import.go"},
		},
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/..."},
			[]string{p + "github.com/avito-tech/go-mutesting/internal/importing/filepath.go", p + "github.com/avito-tech/go-mutesting/internal/importing/import.go"},
		},
	} {
		var opts = &models.Options{}
		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

func TestPackagesWithFilesOfArgs(t *testing.T) {
	//TODO we need normal test with test folder
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
			[]string{"../importing"},
			[]Package{Package{Name: "../importing", Files: []string{"../importing/filepath.go", "../importing/import.go"}}},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing"},
			[]Package{Package{Name: p + "github.com/avito-tech/go-mutesting/internal/importing", Files: []string{p + "github.com/avito-tech/go-mutesting/internal/importing/filepath.go", p + "github.com/avito-tech/go-mutesting/internal/importing/import.go"}}},
		},
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/..."},
			[]Package{
				Package{Name: p + "github.com/avito-tech/go-mutesting/internal/importing", Files: []string{p + "github.com/avito-tech/go-mutesting/internal/importing/filepath.go", p + "github.com/avito-tech/go-mutesting/internal/importing/import.go"}},
			},
		},
	} {
		var opts = &models.Options{}
		got := PackagesWithFilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

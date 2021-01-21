package importing

import (
	"fmt"
	"github.com/avito-tech/go-mutesting/internal/models"
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
			[]string{"./filepath_fixtures/first.go"},
			[]string{"./filepath_fixtures/first.go"},
		},
		// directories
		{
			[]string{"./filepath_fixtures"},
			[]string{"filepath_fixtures/first.go", "filepath_fixtures/second.go", "filepath_fixtures/third.go"},
		},
		{
			[]string{"../importing/filepath_fixtures"},
			[]string{
				"../importing/filepath_fixtures/first.go",
				"../importing/filepath_fixtures/second.go",
				"../importing/filepath_fixtures/third.go",
			},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures"},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/first.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/third.go",
			},
		},
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/..."},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/first.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/third.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second_fixtures_package/fourth.go",
			},
		},
	} {
		var opts = &models.Options{}
		got := FilesOfArgs(test.args, opts)

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
			[]Package{{Name: ".", Files: []string{"filepath.go", "import.go"}}},
		},
		// files
		{
			[]string{"./filepath_fixtures/first.go"},
			[]Package{{Name: "filepath_fixtures", Files: []string{"./filepath_fixtures/first.go"}}},
		},
		// directories
		{
			[]string{"./filepath_fixtures"},
			[]Package{{Name: "filepath_fixtures", Files: []string{
				"filepath_fixtures/first.go",
				"filepath_fixtures/second.go",
				"filepath_fixtures/third.go",
			}}},
		},
		{
			[]string{"../importing/filepath_fixtures"},
			[]Package{{Name: "../importing/filepath_fixtures", Files: []string{
				"../importing/filepath_fixtures/first.go",
				"../importing/filepath_fixtures/second.go",
				"../importing/filepath_fixtures/third.go",
			}}},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures"},
			[]Package{{
				Name: p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures",
				Files: []string{
					p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/first.go",
					p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second.go",
					p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/third.go",
				},
			}},
		},
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/..."},
			[]Package{
				{
					Name: p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures",
					Files: []string{
						p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/first.go",
						p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second.go",
						p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/third.go",
					},
				},
				{
					Name: p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second_fixtures_package",
					Files: []string{
						p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second_fixtures_package/fourth.go",
					},
				},
			},
		},
	} {
		var opts = &models.Options{}
		got := PackagesWithFilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

func TestFilesWithSkipWithoutTests(t *testing.T) {
	p := os.Getenv("GOPATH") + "/src/"

	for _, test := range []struct {
		args   []string
		expect []string
	}{
		// files
		{
			[]string{"./filepath_fixtures/first.go"},
			[]string(nil),
		},
		{
			[]string{"./filepath_fixtures/second.go"},
			[]string{"./filepath_fixtures/second.go"},
		},
		// directories
		{
			[]string{"./filepath_fixtures"},
			[]string{"filepath_fixtures/second.go", "filepath_fixtures/third.go"},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/..."},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/third.go",
			},
		},
	} {
		var opts = &models.Options{}
		opts.Config.SkipFileWithoutTest = true
		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

func TestFilesWithSkipWithBuildTagsTests(t *testing.T) {
	p := os.Getenv("GOPATH") + "/src/"

	for _, test := range []struct {
		args   []string
		expect []string
	}{
		// files
		{
			[]string{"./filepath_fixtures/first.go"},
			[]string(nil),
		},
		{
			[]string{"./filepath_fixtures/third.go"},
			[]string(nil),
		},
		{
			[]string{"./filepath_fixtures/second.go"},
			[]string{"./filepath_fixtures/second.go"},
		},
		// directories
		{
			[]string{"./filepath_fixtures"},
			[]string{"filepath_fixtures/second.go"},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/..."},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepath_fixtures/second.go",
			},
		},
	} {
		var opts = &models.Options{}
		opts.Config.SkipFileWithBuildTag = true
		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

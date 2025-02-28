package importing

import (
	"fmt"
	"os"
	"testing"

	"github.com/avito-tech/go-mutesting/internal/models"

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
			[]string{"./filepathfixtures/first.go"},
			[]string{"./filepathfixtures/first.go"},
		},
		// directories
		{
			[]string{"./filepathfixtures"},
			[]string{"filepathfixtures/first.go", "filepathfixtures/second.go", "filepathfixtures/third.go"},
		},
		{
			[]string{"../importing/filepathfixtures"},
			[]string{
				"../importing/filepathfixtures/first.go",
				"../importing/filepathfixtures/second.go",
				"../importing/filepathfixtures/third.go",
			},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures"},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/first.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/second.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/third.go",
			},
		},
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/first.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/second.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/third.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/secondfixturespackage/fourth.go",
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
			[]string{"./filepathfixtures/first.go"},
			[]Package{{Name: "filepathfixtures", Files: []string{"./filepathfixtures/first.go"}}},
		},
		// directories
		{
			[]string{"./filepathfixtures"},
			[]Package{{Name: "filepathfixtures", Files: []string{
				"filepathfixtures/first.go",
				"filepathfixtures/second.go",
				"filepathfixtures/third.go",
			}}},
		},
		{
			[]string{"../importing/filepathfixtures"},
			[]Package{{Name: "../importing/filepathfixtures", Files: []string{
				"../importing/filepathfixtures/first.go",
				"../importing/filepathfixtures/second.go",
				"../importing/filepathfixtures/third.go",
			}}},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures"},
			[]Package{{
				Name: p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures",
				Files: []string{
					p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/first.go",
					p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/second.go",
					p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/third.go",
				},
			}},
		},
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."},
			[]Package{
				{
					Name: p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures",
					Files: []string{
						p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/first.go",
						p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/second.go",
						p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/third.go",
					},
				},
				{
					Name: p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/secondfixturespackage",
					Files: []string{
						p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/secondfixturespackage/fourth.go",
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
			[]string{"./filepathfixtures/first.go"},
			[]string(nil),
		},
		{
			[]string{"./filepathfixtures/second.go"},
			[]string{"./filepathfixtures/second.go"},
		},
		// directories
		{
			[]string{"./filepathfixtures"},
			[]string{"filepathfixtures/second.go", "filepathfixtures/third.go"},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/second.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/third.go",
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
			[]string{"./filepathfixtures/first.go"},
			[]string(nil),
		},
		{
			[]string{"./filepathfixtures/third.go"},
			[]string(nil),
		},
		{
			[]string{"./filepathfixtures/second.go"},
			[]string{"./filepathfixtures/second.go"},
		},
		// directories
		{
			[]string{"./filepathfixtures"},
			[]string{"filepathfixtures/second.go"},
		},
		// packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/second.go",
			},
		},
	} {
		var opts = &models.Options{}
		opts.Config.SkipFileWithBuildTag = true
		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

func TestFilesWithExcludedDirs(t *testing.T) {
	p := os.Getenv("GOPATH") + "/src/"

	for _, test := range []struct {
		args   []string
		expect []string
		config []string
	}{
		// files
		{
			[]string{"./filepathfixtures/first.go"},
			[]string{"./filepathfixtures/first.go"},
			[]string(nil),
		},
		{
			[]string{"./filepathfixtures/second.go"},
			[]string{"./filepathfixtures/second.go"},
			[]string{"filepathfixtures"},
		},
		{
			[]string{"filepathfixtures/second.go"},
			[]string(nil),
			[]string{"filepathfixtures"},
		},
		{
			[]string{"./filepathfixtures/second.go"},
			[]string(nil),
			[]string{"./filepathfixtures"},
		},
		// directories
		{
			[]string{"./filepathfixtures/..."},
			[]string{
				"filepathfixtures/first.go",
				"filepathfixtures/second.go",
				"filepathfixtures/third.go",
			},
			[]string{"filepathfixtures/secondfixturespackage"},
		},
		{
			[]string{"./filepathfixtures/..."},
			[]string(nil),
			[]string{"filepathfixtures"},
		},
		{
			[]string{"./filepathfixtures"},
			[]string(nil),
			[]string{"filepathfixtures"},
		},
		{
			[]string{"./filepathfixtures"},
			[]string{
				"filepathfixtures/first.go",
				"filepathfixtures/second.go",
				"filepathfixtures/third.go",
			},
			[]string(nil),
		},

		//packages
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/first.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/second.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/third.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/secondfixturespackage/fourth.go",
			},
			[]string{"filepathfixtures"},
		},
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."},
			[]string{
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/first.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/second.go",
				p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/third.go",
			},
			[]string{p + "github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/secondfixturespackage/"},
		},
	} {
		var opts = &models.Options{}
		opts.Config.ExcludeDirs = test.config

		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}
}

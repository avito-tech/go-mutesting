package importing

import (
	"fmt"
	"go/build"
	"path/filepath"
	"testing"

	"github.com/avito-tech/go-mutesting/internal/models"

	"github.com/stretchr/testify/assert"
)

// fixtureDir returns the absolute path to filepathfixtures on this machine.
func fixtureDir(t *testing.T) string {
	t.Helper()
	pkg, err := build.Import("github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures", ".", 0)
	if err != nil {
		t.Skipf("cannot resolve filepathfixtures: %v", err)
	}
	return pkg.Dir
}

// skipIfNoWildcard skips if wildcard package resolution (which walks $GOPATH/src)
// returns nothing — this happens in module-mode checkouts that lack a $GOPATH/src layout.
func skipIfNoWildcard(t *testing.T) {
	t.Helper()
	files := FilesOfArgs([]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."}, &models.Options{})
	if len(files) == 0 {
		t.Skip("wildcard package resolution requires $GOPATH/src layout; skipping in module mode")
	}
}

func TestFilesOfArgs(t *testing.T) {
	dir := fixtureDir(t)
	sub := filepath.Join(dir, "secondfixturespackage")

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
			[]string{"filepathfixtures/fifth.go", "filepathfixtures/first.go", "filepathfixtures/second.go", "filepathfixtures/third.go"},
		},
		{
			[]string{"../importing/filepathfixtures"},
			[]string{
				"../importing/filepathfixtures/fifth.go",
				"../importing/filepathfixtures/first.go",
				"../importing/filepathfixtures/second.go",
				"../importing/filepathfixtures/third.go",
			},
		},
		// single package (non-wildcard)
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures"},
			[]string{
				filepath.Join(dir, "fifth.go"),
				filepath.Join(dir, "first.go"),
				filepath.Join(dir, "second.go"),
				filepath.Join(dir, "third.go"),
			},
		},
	} {
		var opts = &models.Options{}
		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}

	// Wildcard package resolution walks $GOPATH/src and is broken in module mode.
	t.Run("wildcard package", func(t *testing.T) {
		skipIfNoWildcard(t)
		var opts = &models.Options{}
		got := FilesOfArgs([]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."}, opts)
		assert.Equal(t, []string{
			filepath.Join(dir, "fifth.go"),
			filepath.Join(dir, "first.go"),
			filepath.Join(dir, "second.go"),
			filepath.Join(dir, "third.go"),
			filepath.Join(sub, "fourth.go"),
		}, got)
	})
}

func TestPackagesWithFilesOfArgs(t *testing.T) {
	dir := fixtureDir(t)
	sub := filepath.Join(dir, "secondfixturespackage")

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
				"filepathfixtures/fifth.go",
				"filepathfixtures/first.go",
				"filepathfixtures/second.go",
				"filepathfixtures/third.go",
			}}},
		},
		{
			[]string{"../importing/filepathfixtures"},
			[]Package{{Name: "../importing/filepathfixtures", Files: []string{
				"../importing/filepathfixtures/fifth.go",
				"../importing/filepathfixtures/first.go",
				"../importing/filepathfixtures/second.go",
				"../importing/filepathfixtures/third.go",
			}}},
		},
		// single package (non-wildcard)
		{
			[]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures"},
			[]Package{{
				Name: dir,
				Files: []string{
					filepath.Join(dir, "fifth.go"),
					filepath.Join(dir, "first.go"),
					filepath.Join(dir, "second.go"),
					filepath.Join(dir, "third.go"),
				},
			}},
		},
	} {
		var opts = &models.Options{}
		got := PackagesWithFilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}

	// Wildcard package resolution walks $GOPATH/src and is broken in module mode.
	t.Run("wildcard package", func(t *testing.T) {
		skipIfNoWildcard(t)
		var opts = &models.Options{}
		got := PackagesWithFilesOfArgs([]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."}, opts)
		assert.Equal(t, []Package{
			{
				Name: dir,
				Files: []string{
					filepath.Join(dir, "fifth.go"),
					filepath.Join(dir, "first.go"),
					filepath.Join(dir, "second.go"),
					filepath.Join(dir, "third.go"),
				},
			},
			{
				Name:  sub,
				Files: []string{filepath.Join(sub, "fourth.go")},
			},
		}, got)
	})
}

func TestFilesWithSkipWithoutTests(t *testing.T) {
	dir := fixtureDir(t)

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
			[]string{"filepathfixtures/fifth.go", "filepathfixtures/second.go", "filepathfixtures/third.go"},
		},
	} {
		var opts = &models.Options{}
		opts.Config.SkipFileWithoutTest = true
		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}

	// Wildcard package resolution walks $GOPATH/src and is broken in module mode.
	t.Run("wildcard package", func(t *testing.T) {
		skipIfNoWildcard(t)
		var opts = &models.Options{}
		opts.Config.SkipFileWithoutTest = true
		got := FilesOfArgs([]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."}, opts)
		assert.Equal(t, []string{
			filepath.Join(dir, "fifth.go"),
			filepath.Join(dir, "second.go"),
			filepath.Join(dir, "third.go"),
		}, got)
	})
}

func TestFilesWithSkipWithBuildTagsTests(t *testing.T) {
	dir := fixtureDir(t)

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
			[]string{"./filepathfixtures/fifth.go"},
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
	} {
		var opts = &models.Options{}
		opts.Config.SkipFileWithBuildTag = true
		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}

	// Wildcard package resolution walks $GOPATH/src and is broken in module mode.
	t.Run("wildcard package", func(t *testing.T) {
		skipIfNoWildcard(t)
		var opts = &models.Options{}
		opts.Config.SkipFileWithBuildTag = true
		got := FilesOfArgs([]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."}, opts)
		assert.Equal(t, []string{
			filepath.Join(dir, "second.go"),
		}, got)
	})
}

func TestFilesWithExcludedDirs(t *testing.T) {
	dir := fixtureDir(t)
	sub := filepath.Join(dir, "secondfixturespackage")

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
				"filepathfixtures/fifth.go",
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
				"filepathfixtures/fifth.go",
				"filepathfixtures/first.go",
				"filepathfixtures/second.go",
				"filepathfixtures/third.go",
			},
			[]string(nil),
		},
	} {
		var opts = &models.Options{}
		opts.Config.ExcludeDirs = test.config

		got := FilesOfArgs(test.args, opts)

		assert.Equal(t, test.expect, got, fmt.Sprintf("With args: %#v", test.args))
	}

	// Wildcard package resolution walks $GOPATH/src and is broken in module mode.
	t.Run("wildcard package no exclusion", func(t *testing.T) {
		skipIfNoWildcard(t)
		var opts = &models.Options{}
		opts.Config.ExcludeDirs = []string{"filepathfixtures"}
		got := FilesOfArgs([]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."}, opts)
		assert.Equal(t, []string{
			filepath.Join(dir, "fifth.go"),
			filepath.Join(dir, "first.go"),
			filepath.Join(dir, "second.go"),
			filepath.Join(dir, "third.go"),
			filepath.Join(sub, "fourth.go"),
		}, got)
	})

	t.Run("wildcard package exclude subpackage", func(t *testing.T) {
		skipIfNoWildcard(t)
		var opts = &models.Options{}
		opts.Config.ExcludeDirs = []string{sub + "/"}
		got := FilesOfArgs([]string{"github.com/avito-tech/go-mutesting/internal/importing/filepathfixtures/..."}, opts)
		assert.Equal(t, []string{
			filepath.Join(dir, "fifth.go"),
			filepath.Join(dir, "first.go"),
			filepath.Join(dir, "second.go"),
			filepath.Join(dir, "third.go"),
		}, got)
	})
}

package aferoassert

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Errorf(format string, args ...interface{})
}

type tHelper interface {
	Helper()
}

func stat(fs afero.Fs, path string) (os.FileInfo, error) {
	if fs, ok := fs.(afero.Lstater); ok {
		fi, _, err := fs.LstatIfPossible(path)

		return fi, err
	}

	return fs.Stat(path)
}

// Exists checks whether a file or directory exists in the given path. It also fails if there is an error when trying to
// check the file.
func Exists(t TestingT, fs afero.Fs, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if _, err := stat(fs, path); err != nil {
		if os.IsNotExist(err) {
			return assert.Fail(t, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
		}

		return assert.Fail(t, fmt.Sprintf("error when running stat(%q): %s", path, err), msgAndArgs...)
	}

	return true
}

// NoExists checks whether a file does not exist in a given path.
func NoExists(t TestingT, fs afero.Fs, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if _, err := stat(fs, path); err != nil {
		return true
	}

	return assert.Fail(t, fmt.Sprintf("file %q exists", path), msgAndArgs...)
}

// FileExists checks whether a file exists in the given path. It also fails if
// the path points to a directory or there is an error when trying to check the file.
func FileExists(t TestingT, fs afero.Fs, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	info, err := stat(fs, path)
	if err != nil {
		if os.IsNotExist(err) {
			return assert.Fail(t, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
		}

		return assert.Fail(t, fmt.Sprintf("error when running stat(%q): %s", path, err), msgAndArgs...)
	}

	if info.IsDir() {
		return assert.Fail(t, fmt.Sprintf("%q is a directory", path), msgAndArgs...)
	}

	return true
}

// NoFileExists checks whether a file does not exist in a given path. It fails
// if the path points to an existing _file_ only.
func NoFileExists(t TestingT, fs afero.Fs, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	info, err := stat(fs, path)
	if err != nil {
		return true
	}

	if info.IsDir() {
		return true
	}

	return assert.Fail(t, fmt.Sprintf("file %q exists", path), msgAndArgs...)
}

// DirExists checks whether a directory exists in the given path. It also fails
// if the path is a file rather a directory or there is an error checking whether it exists.
func DirExists(t TestingT, fs afero.Fs, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	info, err := stat(fs, path)
	if err != nil {
		if os.IsNotExist(err) {
			return assert.Fail(t, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
		}

		return assert.Fail(t, fmt.Sprintf("error when running stat(%q): %s", path, err), msgAndArgs...)
	}

	if !info.IsDir() {
		return assert.Fail(t, fmt.Sprintf("%q is a file", path), msgAndArgs...)
	}

	return true
}

// NoDirExists checks whether a directory does not exist in the given path.
// It fails if the path points to an existing _directory_ only.
func NoDirExists(t TestingT, fs afero.Fs, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	info, err := stat(fs, path)
	if err != nil {
		return true
	}

	if !info.IsDir() {
		return true
	}

	return assert.Fail(t, fmt.Sprintf("directory %q exists", path), msgAndArgs...)
}

// Perm checks whether a path has the expected permission or not.
func Perm(t TestingT, fs afero.Fs, path string, expected os.FileMode, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	info, err := stat(fs, path)
	if err != nil {
		return assert.Fail(t, fmt.Sprintf("error when running stat(%q): %s", path, err), msgAndArgs...)
	}

	actual := info.Mode() & os.ModePerm

	if expected != actual {
		return assert.Fail(t, fmt.Sprintf("%q permission is 0%o, expected 0%o", path, actual, expected), msgAndArgs...)
	}

	return true
}

// FileContent checks whether a file content is as expected or not.
func FileContent(t TestingT, fs afero.Fs, path string, expected string, msgAndArgs ...interface{}) bool {
	if !FileExists(t, fs, path, msgAndArgs...) {
		return false
	}

	f, err := fs.Open(path)
	if err != nil {
		return assert.Fail(t, fmt.Sprintf("could not open %q: %s", path, err), msgAndArgs...)
	}

	defer f.Close() // nolint: errcheck

	buf := new(bytes.Buffer)

	if _, err := io.Copy(buf, f); err != nil {
		return assert.Fail(t, fmt.Sprintf("could not read %q: %s", path, err), msgAndArgs...)
	}

	return assert.Equal(t, expected, buf.String(), msgAndArgs...)
}

// TreeEqual checks whether a directory is the same as the expectation or not.
func TreeEqual(t TestingT, fs afero.Fs, tree FileTree, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	return assertTree(t, fs, tree, path, true, msgAndArgs...)
}

// YAMLTreeEqual checks whether a directory is the same as the expectation or not.
func YAMLTreeEqual(t TestingT, fs afero.Fs, expected, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	var ft FileTree

	if err := yaml.Unmarshal([]byte(expected), &ft); err != nil {
		return assert.Fail(t, "could not unmarshal expectation", msgAndArgs...)
	}

	return TreeEqual(t, fs, ft, path, msgAndArgs...)
}

// TreeContains checks whether a directory contains a file tree or not.
func TreeContains(t TestingT, fs afero.Fs, tree FileTree, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	return assertTree(t, fs, tree, path, false, msgAndArgs...)
}

// YAMLTreeContains checks whether a directory contains a file tree or not.
func YAMLTreeContains(t TestingT, fs afero.Fs, expected, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	var ft FileTree

	if err := yaml.Unmarshal([]byte(expected), &ft); err != nil {
		return assert.Fail(t, "could not unmarshal expectation", msgAndArgs...)
	}

	return TreeContains(t, fs, ft, path, msgAndArgs...)
}

// nolint: funlen, cyclop
func assertTree(t TestingT, fs afero.Fs, tree FileTree, root string, exhaustive bool, msgAndArgs ...interface{}) bool {
	root = filepath.Clean(root)
	expectations := tree.Flatten("")
	result := true

	fail := func(failureMessage string, args ...interface{}) bool {
		result = false

		return assert.Fail(t, fmt.Sprintf(failureMessage, args...), msgAndArgs...)
	}

	err := afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		expectedPath := strings.TrimPrefix(path, root+string(os.PathSeparator))
		expected, ok := expectations[expectedPath]

		if !ok {
			if exhaustive {
				fail("unexpected file %q", path)
			}

			return nil
		}

		if expected.isDir {
			if !info.IsDir() {
				fail("%q is not a directory", path)

				return nil
			}
		} else if info.IsDir() {
			fail("%q is a directory", path)

			return nil
		}

		if m := expected.Tags.Mode(); m != nil {
			expected := fileModeToString(*m)
			actual := fileModeToString(info.Mode())

			if expected != actual {
				fail("%q mode is %s, expected %s", path, actual, expected)
			}
		}

		if expected := expected.Tags.Perm(); expected != nil {
			actual := info.Mode() & os.ModePerm

			if *expected != actual {
				fail("%q perm is 0%o, expected 0%o", path, actual, *expected)
			}
		}

		delete(expectations, expectedPath)

		return nil
	})
	if err != nil {
		return fail("could not walk through %q: %s", root, err)
	}

	if !result {
		return false
	}

	if len(expectations) == 0 {
		return true
	}

	var sb strings.Builder

	_, _ = sb.WriteString("expected these files in %q but not found:\n")

	for k := range expectations {
		_, _ = fmt.Fprintf(&sb, "- %s\n", k)
	}

	return fail(sb.String(), root)
}

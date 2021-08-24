package aferoassert

import (
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/nhatthm/aferomock"
	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTempSymlinkPath(file string) (string, error) {
	link := file + "_symlink"
	err := os.Symlink(file, link)

	return link, err
}

func cleanUpTempFiles(paths []string) []error {
	var res []error

	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			res = append(res, err)
		}
	}

	return res
}

func TestExists(t *testing.T) {
	osFs := afero.NewOsFs()

	mockT := new(testing.T)
	assert.True(t, Exists(mockT, osFs, "assertions.go"))

	mockT = new(testing.T)
	assert.False(t, Exists(mockT, osFs, "random_file"))

	mockT = new(testing.T)
	assert.True(t, Exists(mockT, osFs, ".github"))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.True(t, Exists(mockT, osFs, link))

	link, err = getTempSymlinkPath("non_existent_file")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)

	mockT = new(testing.T)
	assert.True(t, Exists(mockT, osFs, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestExists_CouldNotStat(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github").
			Return(nil, errors.New("stat error"))
	})(t)

	mockT := new(testing.T)
	assert.False(t, Exists(mockT, fs, ".github"))
}

func TestNoExists(t *testing.T) {
	osFs := afero.NewOsFs()

	mockT := new(testing.T)
	assert.False(t, NoExists(mockT, osFs, "assertions.go"))

	mockT = new(testing.T)
	assert.True(t, NoExists(mockT, osFs, "non_existent_file"))

	mockT = new(testing.T)
	assert.False(t, NoExists(mockT, osFs, ".github"))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.False(t, NoExists(mockT, osFs, link))

	link, err = getTempSymlinkPath("non_existent_file")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.False(t, NoExists(mockT, osFs, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestNoExists_CouldNotStat(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github").
			Return(nil, errors.New("stat error"))
	})(t)

	mockT := new(testing.T)
	assert.True(t, NoExists(mockT, fs, ".github"))
}

func TestFileExists(t *testing.T) {
	osFs := afero.NewOsFs()

	mockT := new(testing.T)
	assert.True(t, FileExists(mockT, osFs, "assertions.go"))

	mockT = new(testing.T)
	assert.False(t, FileExists(mockT, osFs, "random_file"))

	mockT = new(testing.T)
	assert.False(t, FileExists(mockT, osFs, ".github"))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.True(t, FileExists(mockT, osFs, link))

	link, err = getTempSymlinkPath("non_existent_file")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)

	mockT = new(testing.T)
	assert.True(t, FileExists(mockT, osFs, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestFileExists_CouldNotStat(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github").
			Return(nil, errors.New("stat error"))
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileExists(mockT, fs, ".github"))
}

func TestNoFileExists(t *testing.T) {
	osFs := afero.NewOsFs()

	mockT := new(testing.T)
	assert.False(t, NoFileExists(mockT, osFs, "assertions.go"))

	mockT = new(testing.T)
	assert.True(t, NoFileExists(mockT, osFs, "non_existent_file"))

	mockT = new(testing.T)
	assert.True(t, NoFileExists(mockT, osFs, ".github"))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.False(t, NoFileExists(mockT, osFs, link))

	link, err = getTempSymlinkPath("non_existent_file")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.False(t, NoFileExists(mockT, osFs, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestDirExists(t *testing.T) {
	osFs := afero.NewOsFs()

	mockT := new(testing.T)
	assert.False(t, DirExists(mockT, osFs, "assertions.go"))

	mockT = new(testing.T)
	assert.False(t, DirExists(mockT, osFs, "non_existent_dir"))

	mockT = new(testing.T)
	assert.True(t, DirExists(mockT, osFs, ".github"))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.False(t, DirExists(mockT, osFs, link))

	link, err = getTempSymlinkPath("non_existent_dir")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.False(t, DirExists(mockT, osFs, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestDirExists_CouldNotStat(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github").
			Return(nil, errors.New("stat error"))
	})(t)

	mockT := new(testing.T)
	assert.False(t, DirExists(mockT, fs, ".github"))
}

func TestNoDirExists(t *testing.T) {
	osFs := afero.NewOsFs()

	mockT := new(testing.T)
	assert.True(t, NoDirExists(mockT, osFs, "assertions.go"))

	mockT = new(testing.T)
	assert.True(t, NoDirExists(mockT, osFs, "non_existent_dir"))

	mockT = new(testing.T)
	assert.False(t, NoDirExists(mockT, osFs, ".github"))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.True(t, NoDirExists(mockT, osFs, link))

	link, err = getTempSymlinkPath("non_existent_dir")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}

	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assert.True(t, NoDirExists(mockT, osFs, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestPerm(t *testing.T) {
	t.Parallel()

	osFs := afero.NewOsFs()

	mockT := new(testing.T)
	assert.True(t, Perm(mockT, osFs, "assertions.go", 0o644))

	mockT = new(testing.T)
	assert.False(t, Perm(mockT, osFs, "assertions.go", 0o755))
}

func TestPerm_CouldNotStat(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github").
			Return(nil, errors.New("stat error"))
	})(t)

	mockT := new(testing.T)
	assert.False(t, Perm(mockT, fs, ".github", 0o644))
}

func TestFileContent_Success(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.MkdirAll(".github", 0o6444)
	require.NoError(t, err)

	f, err := fs.OpenFile(".github/file.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0o644))
	require.NoError(t, err)

	_, _ = f.WriteString("hello world!") // nolint: errcheck

	mockT := new(testing.T)
	assert.True(t, FileContent(mockT, fs, ".github/file.txt", "hello world!"))

	mockT = new(testing.T)
	assert.False(t, FileContent(mockT, fs, ".github/file.txt", "wrong!"))
}

func TestFileContent_CouldNotStat(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github/file.txt").
			Return(nil, errors.New("stat error"))
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileContent(mockT, fs, ".github/file.txt", "'"))
}

func TestFileContent_FileNotExists(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github/file.txt").
			Return(nil, os.ErrNotExist)
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileContent(mockT, fs, ".github/file.txt", "'"))
}

func TestFileContent_CouldNotOpen(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github/file.txt").
			Return(aferomock.NewFileInfo(func(i *aferomock.FileInfo) {
				i.On("IsDir").Return(false)
			}), nil)

		fs.On("Open", ".github/file.txt").
			Return(nil, errors.New("open error"))
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileContent(mockT, fs, ".github/file.txt", "'"))
}

func TestFileContent_FileIsClosed(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github/file.txt").
			Return(aferomock.NewFileInfo(func(i *aferomock.FileInfo) {
				i.On("IsDir").Return(false)
			}), nil)

		f := mem.NewFileHandle(mem.CreateFile("file.txt"))
		_ = f.Close() // nolint: errcheck

		fs.On("Open", ".github/file.txt").
			Return(f, nil)
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileContent(mockT, fs, ".github/file.txt", "'"))
}

func TestFileContentRegexp_Success(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.MkdirAll(".github", 0o644)
	require.NoError(t, err)

	f, err := fs.OpenFile(".github/file.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0o644))
	require.NoError(t, err)

	_, _ = f.WriteString("hello world!") // nolint: errcheck

	mockT := new(testing.T)
	assert.True(t, FileContentRegexp(mockT, fs, ".github/file.txt", "hello [^!]+!"))

	mockT = new(testing.T)
	assert.True(t, FileContentRegexp(mockT, fs, ".github/file.txt", regexp.MustCompile("hello [^!]+!")))

	mockT = new(testing.T)
	assert.False(t, FileContentRegexp(mockT, fs, ".github/file.txt", "hello [^!]+$"))
}

func TestFileContentRegexp_CouldNotStat(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github/file.txt").
			Return(nil, errors.New("stat error"))
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileContentRegexp(mockT, fs, ".github/file.txt", "'"))
}

func TestFileContentRegexp_FileNotExists(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github/file.txt").
			Return(nil, os.ErrNotExist)
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileContentRegexp(mockT, fs, ".github/file.txt", "'"))
}

func TestFileContentRegexp_CouldNotOpen(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github/file.txt").
			Return(aferomock.NewFileInfo(func(i *aferomock.FileInfo) {
				i.On("IsDir").Return(false)
			}), nil)

		fs.On("Open", ".github/file.txt").
			Return(nil, errors.New("open error"))
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileContentRegexp(mockT, fs, ".github/file.txt", "'"))
}

func TestFileContentRegexp_FileIsClosed(t *testing.T) {
	fs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github/file.txt").
			Return(aferomock.NewFileInfo(func(i *aferomock.FileInfo) {
				i.On("IsDir").Return(false)
			}), nil)

		f := mem.NewFileHandle(mem.CreateFile("file.txt"))
		_ = f.Close() // nolint: errcheck

		fs.On("Open", ".github/file.txt").
			Return(f, nil)
	})(t)

	mockT := new(testing.T)
	assert.False(t, FileContentRegexp(mockT, fs, ".github/file.txt", "'"))
}

func TestTreeEqual_Success(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows 'mode:"Dir"':
    - golangci-lint.yaml
    - test.yaml 'perm:"0644"'
`

	mockT := new(testing.T)
	assert.True(t, YAMLTreeEqual(mockT, osFs, tree, ".github"))
}

func TestTreeEqual_Fail_CouldNotMarshal(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `invalid`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeEqual(mockT, osFs, tree, ".github"))
}

func TestTreeEqual_Fail_CouldNotWalk(t *testing.T) {
	osFs := aferomock.MockFs(func(fs *aferomock.Fs) {
		fs.On("Stat", ".github").
			Return(nil, errors.New("stat error"))
	})(t)

	tree := `- workflows:`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeEqual(mockT, osFs, tree, ".github"))
}

func TestTreeEqual_Fail_MoreFilesThanExpected(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows 'mode:"Dir"':
    - golangci-lint.yaml
`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeEqual(mockT, osFs, tree, ".github"))
}

func TestTreeEqual_Fail_ExpectMoreFiles(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows 'mode:"Dir"':
    - golangci-lint.yaml
    - test.yaml 'perm:"0644"'
    - unknown
`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeEqual(mockT, osFs, tree, ".github"))
}

func TestTreeContains_Success(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows 'mode:"Dir"':
    - golangci-lint.yaml
    - test.yaml 'perm:"0644"'
`

	mockT := new(testing.T)
	assert.True(t, YAMLTreeContains(mockT, osFs, tree, ".github"))
}

func TestTreeContains_Fail_CouldNotMarshal(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `invalid`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeContains(mockT, osFs, tree, ".github"))
}

func TestTreeContains_Fail_ExpectMoreFiles(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows 'mode:"Dir"':
    - golangci-lint.yaml
    - test.yaml 'perm:"0644"'
    - unknown
`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeContains(mockT, osFs, tree, ".github"))
}

func TestTreeContains_Fail_WrongMode(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows 'mode:"Dir|Temporary"':
    - golangci-lint.yaml
    - test.yaml
    - unknown
`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeContains(mockT, osFs, tree, ".github"))
}

func TestTreeContains_Fail_WrongPerm(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows 'mode:"Dir"':
    - golangci-lint.yaml
    - test.yaml 'perm:"0755"'
    - unknown
`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeContains(mockT, osFs, tree, ".github"))
}

func TestTreeContains_Fail_FileIsExpected(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `- workflows`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeContains(mockT, osFs, tree, ".github"))
}

func TestTreeContains_Fail_DirIsExpected(t *testing.T) {
	osFs := afero.NewOsFs()

	tree := `
- workflows:
    - test.yaml:
`

	mockT := new(testing.T)
	assert.False(t, YAMLTreeContains(mockT, osFs, tree, ".github"))
}

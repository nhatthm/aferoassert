package aferoassert

import (
	"errors"
	"os"
	"testing"

	"github.com/nhatthm/aferomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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

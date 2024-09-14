package aferoassert

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"gopkg.in/yaml.v3"
)

var (
	// ErrFileNameEmpty indicates that the file name is empty.
	ErrFileNameEmpty = errors.New("file name is empty")
	// ErrInvalidFileTreeFormat indicates that format of the file tree is invalid and could not be parsed.
	ErrInvalidFileTreeFormat = errors.New("invalid file tree format")
	// ErrInvalidFileMode indicates that the file mode is invalid.
	ErrInvalidFileMode = errors.New("invalid file mode")
)

var (
	tagPattern        = regexp.MustCompile("\\s*'[^`]+'$")
	fileModeSeparator = "|"

	fileModeNames = map[os.FileMode]string{
		os.ModeDir:        "Dir",
		os.ModeAppend:     "Append",
		os.ModeExclusive:  "Exclusive",
		os.ModeTemporary:  "Temporary",
		os.ModeSymlink:    "Symlink",
		os.ModeDevice:     "Device",
		os.ModeNamedPipe:  "NamedPipe",
		os.ModeSocket:     "Socket",
		os.ModeSetuid:     "Setuid",
		os.ModeSetgid:     "Setgid",
		os.ModeCharDevice: "CharDevice",
		os.ModeSticky:     "Sticky",
		os.ModeIrregular:  "Irregular",
	}
)

// FileTree is a map of node.
type FileTree map[string]FileNode

// Flatten converts the file tree to a flat map, key is the path to file.
func (t FileTree) Flatten(root string) map[string]FileNode {
	result := make(map[string]FileNode, len(t))

	for _, n := range t {
		for k, nc := range n.Flatten(root) {
			result[k] = nc
		}
	}

	return result
}

// MarshalYAML satisfies yaml.Marshaler.
func (t FileTree) MarshalYAML() (interface{}, error) { // nolint: unparam
	cnt := len(t)

	if cnt == 0 {
		return map[string]interface{}{}, nil
	}

	raw := make([]FileNode, 0, cnt)

	for _, n := range t {
		raw = append(raw, n)
	}

	sort.Slice(raw, func(i, j int) bool {
		return raw[i].Name < raw[j].Name
	})

	return raw, nil
}

// UnmarshalYAML satisfies yaml.Unmarshaler.
func (t *FileTree) UnmarshalYAML(value *yaml.Node) error {
	var raw []FileNode

	if err := value.Decode(&raw); err != nil {
		return err
	}

	*t = make(map[string]FileNode, len(raw))

	for _, n := range raw {
		(*t)[n.Name] = n
	}

	return nil
}

// FileNode contains needed information for assertions.
type FileNode struct {
	Name     string
	Tags     FileModeTags
	Children FileTree
	IsDir    bool
}

// Flatten converts the file tree to a flat map, key is the path to file.
func (n FileNode) Flatten(root string) map[string]FileNode {
	root = filepath.Join(root, n.Name)

	result := make(map[string]FileNode)
	result[root] = n

	for k, v := range n.Children.Flatten(root) {
		result[k] = v
	}

	return result
}

// MarshalYAML satisfies yaml.Marshaler.
func (n FileNode) MarshalYAML() (interface{}, error) { // nolint: unparam
	var nameBld strings.Builder

	_, _ = nameBld.WriteString(n.Name)

	if len(n.Tags) > 0 {
		_, _ = fmt.Fprintf(&nameBld, " '%s'", n.Tags.String())
	}

	if !n.IsDir {
		return nameBld.String(), nil
	}

	raw := map[string]FileTree{nameBld.String(): n.Children}

	return raw, nil
}

// UnmarshalYAML satisfies yaml.Unmarshaler.
func (n *FileNode) UnmarshalYAML(value *yaml.Node) error {
	// nolint: exhaustive
	switch value.Kind {
	case yaml.ScalarNode:
		r, err := unmarshalFile(value)
		if err != nil {
			return err
		}

		*n = *r

	case yaml.MappingNode:
		r, err := unmarshalFolder(value)
		if err != nil {
			return err
		}

		*n = *r

	default:
		return fmt.Errorf("%w, expected !!str or !!map but got %s at line %d", ErrInvalidFileTreeFormat, value.Tag, value.Line)
	}

	return nil
}

// FileModeTags is a list of tagged file mode.
type FileModeTags map[string]*os.FileMode

// Mode returns file mode.
func (t FileModeTags) Mode() *os.FileMode {
	return t["mode"]
}

// Type returns type file mode.
func (t FileModeTags) Type() *os.FileMode {
	return t["type"]
}

// Perm returns perm file mode.
func (t FileModeTags) Perm() *os.FileMode {
	return t["perm"]
}

// String returns tags in struct tag format.
func (t FileModeTags) String() string {
	tags := &structtag.Tags{}

	if m := t.Mode(); m != nil {
		// nolint: errcheck
		_ = tags.Set(&structtag.Tag{
			Key:  "mode",
			Name: fileModeToString(*m),
		})
	}

	if m := t.Type(); m != nil {
		// nolint: errcheck
		_ = tags.Set(&structtag.Tag{
			Key:  "type",
			Name: fileModeToString(*m),
		})
	}

	if m := t.Perm(); m != nil {
		// nolint: errcheck
		_ = tags.Set(&structtag.Tag{
			Key:  "perm",
			Name: fmt.Sprintf("0%o", *m&os.ModePerm),
		})
	}

	return tags.String()
}

func unmarshalFile(value *yaml.Node) (*FileNode, error) {
	var s string

	if err := value.Decode(&s); err != nil {
		return nil, err
	}

	s = strings.Trim(s, " ")

	if tagPattern.MatchString(s) {
		return unmarshalFileWithTags(value)
	}

	if len(s) == 0 {
		return nil, ErrFileNameEmpty
	}

	return &FileNode{Name: s}, nil
}

func unmarshalFileWithTags(value *yaml.Node) (*FileNode, error) {
	rawTags := tagPattern.FindString(value.Value)
	fileName := strings.TrimSuffix(value.Value, rawTags)

	if len(fileName) == 0 {
		return nil, ErrFileNameEmpty
	}

	tags, err := unmarshalTags(value, prepareTagsString(rawTags))
	if err != nil {
		return nil, err
	}

	n := &FileNode{
		Name:     fileName,
		Tags:     *tags,
		Children: nil,
	}

	return n, nil
}

func unmarshalTags(node *yaml.Node, s string) (*FileModeTags, error) {
	tags, err := structtag.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("%w at line %d", err, node.Line)
	}

	t := make(FileModeTags, tags.Len())

	for _, tag := range tags.Tags() {
		value, err := parseTag(tag.Name)
		if err != nil {
			return nil, fmt.Errorf("%w in %q tag at line %d", ErrInvalidFileMode, tag.Key, node.Line)
		}

		t[tag.Key] = value
	}

	return &t, nil
}

func unmarshalFolder(value *yaml.Node) (*FileNode, error) {
	if len(value.Content) != 2 { //nolint: mnd
		return nil, ErrInvalidFileTreeFormat
	}

	d, err := unmarshalFile(value.Content[0])
	if err != nil {
		return nil, err
	}

	var dt FileTree

	if err := value.Content[1].Decode(&dt); err != nil {
		return nil, err
	}

	d.Children = dt
	d.IsDir = true

	return d, nil
}

func prepareTagsString(s string) string {
	return strings.Trim(s, " `'")
}

func parseTag(tag string) (*os.FileMode, error) {
	base := 10
	if strings.HasPrefix(tag, "0") {
		base = 8
	}

	mode, err := strconv.ParseUint(tag, base, 32)
	if err == nil {
		return FileModeFromUint64(mode), nil
	}

	var result os.FileMode

	for _, s := range strings.Split(tag, fileModeSeparator) {
		m, err := fileModeFromString(s)
		if err != nil {
			return nil, err
		}

		result |= *m
	}

	return &result, nil
}

func fileModeFromString(s string) (*os.FileMode, error) {
	for mode, name := range fileModeNames {
		if name == s {
			return &mode, nil
		}
	}

	return nil, ErrInvalidFileMode
}

func fileModeToString(mode os.FileMode) string {
	result := make([]string, 0)

	for m, name := range fileModeNames {
		if mode&m != 0 {
			result = append(result, name)
		}
	}

	sort.Strings(result)

	return strings.Join(result, fileModeSeparator)
}

// FileModePtr returns pointer to file mode.
func FileModePtr(mode os.FileMode) *os.FileMode {
	return &mode
}

// FileModeFromUint64 returns *os.FileMode from an uint64.
func FileModeFromUint64(mode uint64) *os.FileMode {
	result := os.FileMode(mode) //nolint: gosec

	return &result
}

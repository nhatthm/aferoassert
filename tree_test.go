package aferoassert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestFileTree_Flatten(t *testing.T) {
	t.Parallel()

	text := `
- file 1
- folder 2:
    - file 2 'perm:"0755"'
    - folder 3 'mode:"Dir|Sticky" perm:"0644"':
        - file 3
    - file 4
    - folder 4:
`

	var ft FileTree

	err := yaml.Unmarshal([]byte(text), &ft)
	require.NoError(t, err)

	expected := map[string]FileNode{
		"file 1": {Name: "file 1"},
		"folder 2": {
			Name:  "folder 2",
			isDir: true,
			Children: FileTree{
				"file 2": {
					Name: "file 2",
					Tags: FileModeTags{
						"perm": FileModeFromUint64(0o755),
					},
				},
				"folder 3": {
					Name:  "folder 3",
					isDir: true,
					Tags: FileModeTags{
						"mode": FileModePtr(os.ModeDir | os.ModeSticky),
						"perm": FileModeFromUint64(0o644),
					},
					Children: FileTree{
						"file 3": {Name: "file 3"},
					},
				},
				"file 4":   {Name: "file 4"},
				"folder 4": {Name: "folder 4", isDir: true},
			},
		},
		"folder 2/file 2": {
			Name: "file 2",
			Tags: FileModeTags{
				"perm": FileModeFromUint64(0o755),
			},
		},
		"folder 2/folder 3": {
			Name:  "folder 3",
			isDir: true,
			Tags: FileModeTags{
				"mode": FileModePtr(os.ModeDir | os.ModeSticky),
				"perm": FileModeFromUint64(0o644),
			},
			Children: FileTree{
				"file 3": {Name: "file 3"},
			},
		},
		"folder 2/folder 3/file 3": {Name: "file 3"},
		"folder 2/file 4":          {Name: "file 4"},
		"folder 2/folder 4":        {Name: "folder 4", isDir: true},
	}

	assert.Equal(t, expected, ft.Flatten(""))
}

func TestNode_Serde(t *testing.T) {
	t.Parallel()

	text := `
- file 1
- folder 2:
    - file 2 'perm:"0755"'
    - folder 3 'mode:"Dir|Sticky" type:"Dir" perm:"0644"':
        - file 3
    - file 4
    - folder 4 'mode:"Dir|Temporary"':
`

	var ft FileTree

	err := yaml.Unmarshal([]byte(text), &ft)
	require.NoError(t, err)

	result, err := yaml.Marshal(ft)
	require.NoError(t, err)

	expected := `- file 1
- folder 2:
    - file 2 'perm:"0755"'
    - file 4
    - folder 3 'mode:"Dir|Sticky" type:"Dir" perm:"0644"':
        - file 3
    - folder 4 'mode:"Dir|Temporary"': {}
`

	assert.Equal(t, expected, string(result))
}

func TestNode_UnmarshalYAML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		text           string
		expectedResult FileTree
		expectedError  string
	}{
		{
			scenario: "invalid file tree at root",
			text: `
file 1
file 2:
    - file 2
`,
			expectedError: "yaml: line 3: mapping values are not allowed in this context",
		},
		{
			scenario:      "empty file name",
			text:          `- ""`,
			expectedError: `file name is empty`,
		},
		{
			scenario:      "empty file name with tags",
			text:          `- "'mode:\"Dir\"'"`,
			expectedError: `file name is empty`,
		},
		{
			scenario:      "invalid file node",
			text:          `- []`,
			expectedError: `invalid file tree format, expected !!str or !!map but got !!seq at line 1`,
		},
		{
			scenario:      "malformed tag",
			text:          `- file 1 'type:"Unknown'`,
			expectedError: `bad syntax for struct tag value at line 1`,
		},
		{
			scenario:      "invalid file mode at first level (file)",
			text:          `- file 1 'type:"Unknown"'`,
			expectedError: `invalid file mode in "type" tag at line 1`,
		},
		{
			scenario:      "invalid file mode at first level (directory)",
			text:          `- folder 1 'type:"Unknown"':`,
			expectedError: `invalid file mode in "type" tag at line 1`,
		},
		{
			scenario:      "invalid file name in directory",
			text:          `- {}:`,
			expectedError: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!map into string",
		},
		{
			scenario: "invalid directory format",
			text: `
- folder name:
  another property is invalid:
`,
			expectedError: "invalid file tree format",
		},
		{
			scenario: "invalid file mode at second level (file)",
			text: `
- folder 1:
    - file 2 'type:"Unknown"'
`,
			expectedError: `invalid file mode in "type" tag at line 3`,
		},
		{
			scenario: "invalid file mode at second level (directory)",
			text: `
- folder 1:
    - folder 2 'type:"Unknown"':
        - file 1
`,
			expectedError: `invalid file mode in "type" tag at line 3`,
		},
		{
			scenario: "valid with tags",
			text: `
- file 1
- folder 2:
    - file 2 'perm:"0755"'
    - folder 3 'type:"Dir|Sticky" perm:"0644"':
        - file 3
    - file 4
    - folder 4:
`,
			expectedResult: FileTree{
				"file 1": {Name: "file 1"},
				"folder 2": {
					Name:  "folder 2",
					isDir: true,
					Children: FileTree{
						"file 2": {
							Name: "file 2",
							Tags: FileModeTags{
								"perm": FileModeFromUint64(0o755),
							},
						},
						"folder 3": {
							Name:  "folder 3",
							isDir: true,
							Tags: FileModeTags{
								"type": FileModePtr(os.ModeDir | os.ModeSticky),
								"perm": FileModeFromUint64(0o644),
							},
							Children: FileTree{
								"file 3": {Name: "file 3"},
							},
						},
						"file 4":   {Name: "file 4"},
						"folder 4": {Name: "folder 4", isDir: true},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			var ft FileTree

			err := yaml.Unmarshal([]byte(tc.text), &ft)

			assert.Equal(t, tc.expectedResult, ft)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

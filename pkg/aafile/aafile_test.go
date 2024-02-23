package aafile

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestAAFile_Match(t *testing.T) {
	testcases := []struct {
		name     string
		raw      string
		path     string
		expected AAFile
	}{
		{
			name:     "empty",
			raw:      "",
			path:     "path/to/file",
			expected: nil,
		},
		{
			name:     "equal path",
			raw:      "path/to/file @user1 @user2",
			path:     "path/to/file",
			expected: AAFile{{Path: "path/to/file", Users: []string{"user1", "user2"}}},
		},
		{
			name:     "equal path with wildcard",
			raw:      "path/to/* @user1 @user2",
			path:     "path/to/file",
			expected: AAFile{{Path: "path/to/*", Users: []string{"user1", "user2"}}},
		},
		{
			name:     "equal with wildcard",
			raw:      "* @user1 @user2",
			path:     "path/to/file",
			expected: AAFile{{Path: "*", Users: []string{"user1", "user2"}}},
		},
		{
			name:     "not equal path",
			raw:      "path/to/file @user1 @user2",
			path:     "path/to/another/file",
			expected: nil,
		},
		{
			name:     "not equal path with wildcard",
			raw:      "path/to/another/* @user1 @user2",
			path:     "path/to/file",
			expected: nil,
		},
		{
			name:     "not equal with wildcard",
			raw:      "another/* @user1 @user2",
			path:     "path/to/file",
			expected: nil,
		},
		{
			name: "match multiple",
			raw: `path/to/* @user1 @user2
		match/* @user3 @user4
		path/to/**/file @user3 @user4`,
			path: "path/to/file",
			expected: AAFile{
				{Path: "path/to/*", Users: []string{"user1", "user2"}},
				{Path: "path/to/**/file", Users: []string{"user3", "user4"}},
			},
		},
		{
			name:     "match wildcard file",
			raw:      `path/to/*.yaml @user1`,
			path:     "path/to/file.yaml",
			expected: AAFile{{Path: "path/to/*.yaml", Users: []string{"user1"}}},
		},
		{
			name:     "match folder",
			raw:      `path/to @user1`,
			path:     "path/to/file",
			expected: AAFile{{Path: "path/to", Users: []string{"user1"}}},
		},
		{
			name:     "match folder wildcard",
			raw:      `path/**/test @user1`,
			path:     "path/to/test/file",
			expected: AAFile{{Path: "path/**/test", Users: []string{"user1"}}},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			aafile, _, err := ParseRaw([]byte(tt.raw))
			require.NoError(t, err)
			matched := aafile.Match(tt.path)
			if diff := cmp.Diff(tt.expected, matched, cmpopts.EquateEmpty(), cmpopts.IgnoreUnexported(Rule{})); diff != "" {
				t.Errorf("Match() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

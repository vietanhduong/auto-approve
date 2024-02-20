package aafile

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func Test_PraseRaw(t *testing.T) {
	testcases := []struct {
		name     string
		raw      string
		expected AAFile
		warnings []string
	}{
		{
			name:     "empty",
			raw:      "",
			expected: nil,
			warnings: nil,
		},
		{
			name:     "single line",
			raw:      "path/to/file @user1 @user2",
			expected: AAFile{{Path: "path/to/file", Users: []string{"user1", "user2"}}},
			warnings: nil,
		},
		{
			name: "multi lines",
			raw: `path/to/file @user1 @user2
path/to/another/file @user3 @user4`,
			expected: AAFile{
				{Path: "path/to/file", Users: []string{"user1", "user2"}},
				{Path: "path/to/another/file", Users: []string{"user3", "user4"}},
			},
		},
		{
			name: "empty line",
			raw: `path/to/file @user1 @user2


`,
			expected: AAFile{{Path: "path/to/file", Users: []string{"user1", "user2"}}},
		},
		{
			name: "comment line",
			raw: `# this is a comment
path/to/file             @user1 @user2
# this is another comment
path/to/another/file     @user3 @user4`,
			expected: AAFile{
				{Path: "path/to/file", Users: []string{"user1", "user2"}},
				{Path: "path/to/another/file", Users: []string{"user3", "user4"}},
			},
		},
		{
			name: "invalid user",
			raw: `
path/to/file @user1 @user2
path/to/another/file @user3 @user4 invalid`,
			expected: AAFile{
				{Path: "path/to/file", Users: []string{"user1", "user2"}},
				{Path: "path/to/another/file", Users: []string{"user3", "user4"}},
			},
			warnings: []string{"line 3: invalid user invalid"},
		},
		{
			name: "no user specified",
			raw: `
path/to/file @user1 @user2
path/to/another/file`,
			expected: AAFile{
				{Path: "path/to/file", Users: []string{"user1", "user2"}},
				{Path: "path/to/another/file", Users: []string{"user1", "user2"}},
			},
		},
		{
			name: "no user specified at top most record",
			raw: `path/to/file
path/to/another/file @user3 @user4
path/to/f2`,
			expected: AAFile{
				{Path: "path/to/file"},
				{Path: "path/to/another/file", Users: []string{"user3", "user4"}},
				{Path: "path/to/f2", Users: []string{"user3", "user4"}},
			},
			warnings: []string{"line 1: no user specified"},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			aafile, warnings, err := ParseRaw([]byte(tt.raw))
			require.NoError(t, err)
			if diff := cmp.Diff(tt.warnings, warnings, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("unexpected warnings (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.expected, aafile, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("unexpected aafile (-want +got):\n%s", diff)
			}
		})
	}
}

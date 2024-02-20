package aafile

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

var spaceRegex = regexp.MustCompile(`\s+`)

func ParseRaw(raw []byte) (AAFile, []string, error) {
	var ret AAFile
	var warnings []string
	var lc int
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	for scanner.Scan() {
		lc++
		// AAFile uses the same format with the GitHub's CODEOWNERS
		// <path> @<user> @<user> ...
		line := strings.TrimSpace(scanner.Text())
		// replace multi-spaces to single space in-case padding cols
		line = spaceRegex.ReplaceAllString(line, " ")
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, " ")
		r := Rule{Path: parts[0]}
		if len(parts) > 1 {
			r.Users = lo.Filter(lo.Map(parts[1:], func(s string, _ int) string {
				if s[0] != '@' {
					warnings = append(warnings, fmt.Sprintf("line %d: invalid user %s", lc, s))
					return ""
				}
				return s[1:]
			}), func(s string, _ int) bool {
				return s != ""
			})
		}

		if len(r.Users) == 0 {
			if len(ret) == 0 { // top most record
				warnings = append(warnings, fmt.Sprintf("line %d: no user specified", lc))
			} else {
				r.Users = ret[len(ret)-1].Users
			}
		}
		ret = append(ret, r)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scanner: %w", err)
	}
	return ret, warnings, nil
}

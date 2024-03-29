package aafile

import (
	"strings"

	"github.com/samber/lo"
)

type Rule struct {
	Path    string
	pattern pattern
	Users   []string
}

type AAFile []Rule

func (f AAFile) Match(path string) AAFile {
	records := lo.Filter(f, func(r Rule, _ int) bool {
		if r.Path == path || r.Path == "*" || strings.HasPrefix(path, r.Path) {
			return true
		}
		return r.pattern.match(path)
		// return ignore.Match(r.Path, path)
	})
	return records
}

func (f AAFile) MatchUser(user string) AAFile {
	return lo.Filter(f, func(r Rule, _ int) bool {
		return lo.Contains(r.Users, user)
	})
}

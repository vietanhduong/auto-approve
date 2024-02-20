package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/vietanhduong/auto-approve/pkg/aafile"
	"github.com/vietanhduong/auto-approve/pkg/logging"
)

var aaFileLocations = []string{
	"AUTOAPPROVE",
	".github/AUTOAPPROVE",
}

func parseAAFile(path string) (aafile.AAFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}

	ret, warnings, err := aafile.ParseRaw(raw)
	if err != nil {
		return nil, fmt.Errorf("parse file %s: %w", path, err)
	}
	for _, w := range warnings {
		logging.Warning(w)
	}
	return ret, nil
}

func discoveryAaFile(repoRoot string) string {
	paths := lo.Filter(aaFileLocations, func(loc string, _ int) bool {
		logging.Debugf("Tried AAFile path: %s", loc)
		_, err := os.Stat(filepath.Join(repoRoot, loc))
		return err == nil
	})
	if len(paths) == 0 {
		return ""
	}
	return filepath.Join(repoRoot, paths[0])
}

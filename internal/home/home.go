package home

import (
	"os"
	"path/filepath"
	"strings"
)

var homedir, _ = os.UserHomeDir()

// Short replaces the actual home path from [Dir] with `~`.
func Short(p string) string {
	if homedir == "" || !strings.HasPrefix(p, homedir) {
		return p
	}
	return filepath.Join("~", strings.TrimPrefix(p, homedir))
}

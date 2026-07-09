package shared

import (
	"fmt"
	"path/filepath"
)

func formatRelativePath(file string, line int) string {
	dir, filename := filepath.Split(file)

	parentDir := filepath.Base(dir)

	grandparentDir := filepath.Base(filepath.Dir(filepath.Clean(dir)))

	if grandparentDir != "." && grandparentDir != "/" {
		return fmt.Sprintf("%s/%s/%s:%d", grandparentDir, parentDir, filename, line)
	}

	return fmt.Sprintf("%s/%s:%d", parentDir, filename, line)
}

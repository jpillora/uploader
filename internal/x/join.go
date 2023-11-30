package x

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Join(overwrite bool, parts ...string) string {
	if overwrite {
		return filepath.Join(parts...)
	}
	return JoinUnique(parts...)
}

func JoinUnique(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	path := filepath.Join(parts...)
	count := 1
	for {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(path)
		prefix := strings.TrimSuffix(path, ext)
		if count > 1 {
			prefix = strings.TrimSuffix(prefix, fmt.Sprintf("-%d", count))
		}
		count++
		path = fmt.Sprintf("%s-%d%s", prefix, count, ext)
	}
	return path
}

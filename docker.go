package isisolated

import (
	"os"
	"path/filepath"
)

// IsInDocker returns current process works
func IsInDocker() bool {
	_, err := os.Stat(filepath.Join("/", ".dockerenv"))
	return err == nil
}

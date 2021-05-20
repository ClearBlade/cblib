package fsutil

import (
	"fmt"
	"os"
	"path"
)

func EnsureDirectory(dir string) error {
	stat, err := os.Stat(path.Clean(dir))
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	return nil
}

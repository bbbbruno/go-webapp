package backup

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
)

func DirHash(path string) (string, error) {
	hash := md5.New()

	if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if _, err := io.WriteString(hash, path); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(hash, "%v", info.IsDir()); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(hash, "%v", info.ModTime()); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(hash, "%v", info.Mode()); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(hash, "%v", info.Name()); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(hash, "%v", info.Size()); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

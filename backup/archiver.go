package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Archiver interface {
	DestFmt(i int64) string
	Archive(src, dest string) error
}

type zipper struct{}

var ZIP Archiver = (*zipper)(nil)

func (z *zipper) Archive(src string, dest string) (err error) {
	if err = os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	w := zip.NewWriter(out)
	defer func() {
		cerr := w.Close()
		if err == nil {
			err = cerr
		}
	}()

	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		f, err := w.Create(path)
		if err != nil {
			return err
		}
		if _, err = io.Copy(f, in); err != nil {
			return err
		}

		return nil
	})
}

func (z *zipper) DestFmt(i int64) string {
	return fmt.Sprintf("%d.zip", i)
}

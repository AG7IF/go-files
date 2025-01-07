package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type File struct {
	dir  string
	base string
	ext  string
}

func NewFile(path string) (File, error) {
	dir, base, ext, err := DecomposePath(path)
	if err != nil {
		return File{}, errors.WithStack(err)
	}

	return File{
		dir:  dir,
		base: base,
		ext:  ext,
	}, nil
}

func DecomposePath(path string) (string, string, string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", "", "", errors.WithStack(err)
	}

	dir := filepath.Dir(abs)

	base := strings.TrimSuffix(filepath.Base(abs), filepath.Ext(abs))

	ext := filepath.Ext(abs)

	return dir, base, ext, nil
}

func (f File) Empty() bool {
	return f.dir == "" && f.base == "" && f.ext == ""
}

func (f File) Dir() string {
	return f.dir
}

func (f File) Base() string {
	return f.base
}

func (f File) Ext() string {
	return f.ext
}

func (f File) Name() string {
	return fmt.Sprintf("%s%s", f.base, f.ext)
}

func (f File) FullPath() string {
	return filepath.Join(f.dir, fmt.Sprintf("%s%s", f.base, f.ext))
}

func (f File) Stat() (os.FileInfo, error) {
	return os.Stat(f.FullPath())
}

func (f File) Create() (*os.File, error) {
	return os.Create(f.FullPath())
}

func (f File) Open() (*os.File, error) {
	return os.Open(f.FullPath())
}

func (f File) Remove() error {
	return os.Remove(f.FullPath())
}

func (f File) Copy(destDir string) (destFile File, err error) {
	destPath := filepath.Join(destDir, f.Name())
	destFile, err = NewFile(destPath)
	if err != nil {
		destFile = File{}
		err = errors.WithStack(err)
		return
	}

	dest, err := destFile.Create()
	if err != nil {
		destFile = File{}
		err = errors.WithStack(err)
		return
	}
	defer func(dest *os.File) {
		if cerr := dest.Close(); cerr != nil && err == nil {
			err = errors.WithStack(cerr)
		}
	}(dest)

	src, err := f.Open()
	if err != nil {
		destFile = File{}
		err = errors.WithStack(err)
		return
	}
	defer func(src *os.File) {
		if cerr := src.Close(); cerr != nil && err == nil {
			err = errors.WithStack(cerr)
		}
	}(src)

	_, err = io.Copy(dest, src)
	if err != nil {
		destFile = File{}
		err = errors.WithStack(err)
		return
	}

	return
}

func (f File) CopyAndRename(destDir string, newName string) (destFile File, err error) {
	destPath := filepath.Join(destDir, newName)
	destFile, err = NewFile(destPath)
	if err != nil {
		destFile = File{}
		err = errors.WithStack(err)
		return
	}

	dest, err := destFile.Create()
	if err != nil {
		destFile = File{}
		err = errors.WithStack(err)
		return
	}
	defer func(dest *os.File) {
		if cerr := dest.Close(); cerr != nil && err == nil {
			err = errors.WithStack(cerr)
		}
	}(dest)

	src, err := f.Open()
	if err != nil {
		destFile = File{}
		err = errors.WithStack(err)
		return
	}
	defer func(src *os.File) {
		if cerr := src.Close(); cerr != nil && err == nil {
			err = errors.WithStack(cerr)
		}
	}(src)

	_, err = io.Copy(dest, src)
	if err != nil {
		destFile = File{}
		err = errors.WithStack(err)
		return
	}

	return
}

func (f File) Move(destDir string) (File, error) {
	destFile, err := f.Copy(destDir)
	if err != nil {
		return File{}, errors.WithStack(err)
	}

	err = f.Remove()
	if err != nil {
		return File{}, errors.WithStack(err)
	}

	return destFile, nil
}

func (f File) MoveAndRename(destDir string, newName string) (File, error) {
	destFile, err := f.CopyAndRename(destDir, newName)
	if err != nil {
		return File{}, errors.WithStack(err)
	}

	err = f.Remove()
	if err != nil {
		return File{}, errors.WithStack(err)
	}

	return destFile, nil
}

func (f File) ReadFile() ([]byte, error) {
	return os.ReadFile(f.FullPath())
}

func (f File) WriteBytes(content []byte) (err error) {
	out, err := f.Create()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func(out *os.File) {
		if cerr := out.Close(); cerr != nil && err == nil {
			err = errors.WithStack(cerr)
		}
	}(out)

	_, err = out.Write(content)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = out.Sync()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (f File) WriteString(content string) (err error) {
	out, err := f.Create()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func(out *os.File) {
		if cerr := out.Close(); cerr != nil && err == nil {
			err = errors.WithStack(cerr)
		}
	}(out)

	_, err = out.WriteString(content)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = out.Sync()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

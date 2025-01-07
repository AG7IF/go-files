package files

import (
	"fmt"
	"io/fs"
	"os"
	"sort"

	"github.com/pkg/errors"
)

type Directory struct {
	path           string
	subdirectories []*Directory
	files          []File
}

func NewDirectory(path string) *Directory {
	return &Directory{
		path: path,
	}
}

func (d *Directory) sortSubdirectories() {
	sort.Slice(d.subdirectories, func(i, j int) bool {
		return d.subdirectories[i].Path() < d.subdirectories[j].Path()
	})
}

func (d *Directory) sortFiles() {
	sort.Slice(d.files, func(i, j int) bool {
		return d.files[i].Name() < d.files[j].Name()
	})
}

func (d *Directory) Populate() error {
	var subdirectories []*Directory
	var files []File

	fileSystem := os.DirFS(d.path)

	err := fs.WalkDir(fileSystem, ".", func(p string, d fs.DirEntry, e error) error {
		if d.Name() == "." {
			return errors.WithStack(nil)
		}

		if d.IsDir() {
			directory := NewDirectory(d.Name())
			err := directory.Populate()
			if err != nil {
				return errors.WithStack(err)
			}

			subdirectories = append(subdirectories, directory)
			return nil
		}

		file, err := NewFile(d.Name())
		if err != nil {
			return errors.WithStack(err)
		}

		files = append(files, file)
		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	d.subdirectories = subdirectories
	d.sortSubdirectories()

	d.files = files
	d.sortFiles()

	return nil
}

func (d *Directory) Path() string {
	return d.path
}

func (d *Directory) Subdirectories() []*Directory {
	return d.subdirectories
}

func (d *Directory) Files() []File {
	return d.files
}

func (d *Directory) FilterByExt(ext string) []File {
	var filtered []File

	for _, file := range d.files {
		if file.Ext() == fmt.Sprintf(".%s", ext) {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

func (d *Directory) CreateFile(name string) (File, error) {
	f, err := NewFile(fmt.Sprintf("%s/%s", d.path, name))
	if err != nil {
		return File{}, errors.WithStack(err)
	}

	d.files = append(d.files, f)
	d.sortFiles()

	return f, nil
}

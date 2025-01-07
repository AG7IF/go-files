package files

import (
	"io/fs"
	"os"
	"sort"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Directory struct {
	path           string
	subdirectories []Directory
	files          []File
	logger         *zerolog.Logger
}

func NewDirectory(path string, logger *zerolog.Logger) (Directory, error) {
	var subdirectories []Directory
	var files []File

	fileSystem := os.DirFS(path)

	err := fs.WalkDir(fileSystem, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		if d.IsDir() {
			directory, err := NewDirectory(d.Name(), logger)
			if err != nil {
				return errors.WithStack(err)
			}

			subdirectories = append(subdirectories, directory)
			return nil
		}

		file, err := NewFile(d.Name(), logger)
		if err != nil {
			return errors.WithStack(err)
		}

		files = append(files, file)
		return nil
	})
	if err != nil {
		return Directory{}, err
	}

	sort.Slice(subdirectories, func(i, j int) bool {
		return subdirectories[i].Path() < subdirectories[j].Path()
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	return Directory{
		path:           path,
		subdirectories: subdirectories,
		files:          files,
		logger:         logger,
	}, nil
}

func (d Directory) Filter(ext string) []File {
	var filtered []File

	for _, file := range d.files {
		if file.Ext() == ext {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

func (d Directory) Path() string {
	return d.path
}

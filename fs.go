package klaro

import (
	"io/fs"
)

type PrefixFS struct {
	fs     fs.FS
	prefix string
}

func (f *PrefixFS) Open(name string) (fs.File, error) {
	return f.fs.Open(name[len(f.prefix):])
}

type MultiFS struct {
	fileSystems []fs.FS
}

func (m *MultiFS) Open(name string) (fs.File, error) {
	var err error
	var file fs.File
	for _, fileSystem := range m.fileSystems {
		if file, err = fileSystem.Open(name); err == nil {
			return file, nil
		}
	}
	return nil, err
}

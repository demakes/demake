package klaro

type FileBackend struct {
	opts FileBackendOptions
}

type FileBackendOptions struct {
	Path string
}

func MakeFileBackend(opts FileBackendOptions) *FileBackend {
	return &FileBackend{
		opts: opts,
	}
}

package klaro

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

type Server struct {
	opts       Options
	fs         fs.FS
	server     *http.Server
	fileServer http.Handler
	backend    Backend
}

type Options struct {
	StaticPrefix string
}

func MakeServer(opts Options, backend Backend) *Server {

	fs := &PrefixFS{
		fs:     &MultiFS{},
		prefix: opts.StaticPrefix,
	}

	return &Server{
		fs:         fs,
		opts:       opts,
		backend:    backend,
		fileServer: http.FileServer(http.FS(fs)),
		server: &http.Server{
			Addr: ":8001",
		},
	}
}

func computeETag(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf(`"%s"`, hex.EncodeToString(hash[:]))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, s.opts.StaticPrefix) {

		filePath := filepath.Join(".", r.URL.Path)

		// Check if the file exists and is not a directory
		file, err := s.fs.Open(filePath)

		if err == nil {
			defer file.Close()
			fileInfo, err := file.Stat()
			if err == nil && !fileInfo.IsDir() {
				// Read file contents to compute the ETag
				fileContents, err := ioutil.ReadAll(file)
				if err == nil {
					etag := computeETag(fileContents)
					w.Header().Set("ETag", etag)
					w.Header().Set("Cache-Control", "private, max-age=0, must-revalidate")
				}
			}
		}

		// If the ETag in the request matches the computed ETag, return 304 Not Modified
		if r.Header.Get("If-None-Match") == w.Header().Get("ETag") {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		s.fileServer.ServeHTTP(w, r)
		return
	}

	w.Header().Add("content-type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte("test"))
}

func (s *Server) Start() error {
	s.server.Handler = s
	go s.server.ListenAndServe()
	return nil
}

func (s *Server) Stop() {
	// to do: implement stop
}

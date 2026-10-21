package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

// Only what the pages need is embedded; nothing else is reachable over HTTP.
//
//go:embed static
var staticFS embed.FS

// staticHandler serves files from the embedded static tree. It never lists
// directories. Files are revalidated on every request so a new release is
// picked up at once.
func staticHandler() (http.Handler, error) {
	root, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, err
	}

	started := time.Now()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

		info, err := fs.Stat(root, name)
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}

		file, err := root.Open(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer func() { _ = file.Close() }()

		seeker, ok := file.(interface {
			Read([]byte) (int, error)
			Seek(int64, int) (int64, error)
		})
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Cache-Control", "no-cache")
		http.ServeContent(w, r, name, started, seeker)
	}), nil
}

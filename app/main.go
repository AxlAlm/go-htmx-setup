package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	fileServer := http.FileServer(http.FS(fixedFileSystem{staticFiles}))
	fmt.Println(fileServer)
	http.Handle("/static", ContentTypeWrapper(fileServer))
	fmt.Println("Serving at http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

type fixedFileSystem struct {
	fs embed.FS
}

func (f fixedFileSystem) Open(path string) (fs.File, error) {
	// Strip the leading "/" added by http.FileServer
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	fmt.Println(path)
	// Default case for serving files
	return f.fs.Open("static/" + path)
}

// ContentTypeWrapper wraps the FileServer handler to serve the correct content type.
func ContentTypeWrapper(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check file extension to set Content-Type
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		h.ServeHTTP(w, r)
	}
}

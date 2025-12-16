package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Serve(port, outputDir string) {
	// Custom Handler to support Pretty URLs (No trailing slash)
	fileServer := http.FileServer(http.Dir(outputDir))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Redirect trailing slashes to slash-less path (except root)
		// e.g. /my-note/ -> /my-note
		if path != "/" && strings.HasSuffix(path, "/") {
			newPath := strings.TrimSuffix(path, "/")
			http.Redirect(w, r, newPath, http.StatusMovedPermanently)
			return
		}

		// Serve directory index as the file
		// If path has no extension and matches a directory, serve index.html directly
		if path != "/" && filepath.Ext(path) == "" {
			localPath := filepath.Join(outputDir, path)
			if info, err := os.Stat(localPath); err == nil && info.IsDir() {
				indexPath := filepath.Join(localPath, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					http.ServeFile(w, r, indexPath)
					return
				}
			}
		}

		// Fallback to standard file server
		fileServer.ServeHTTP(w, r)
	})

	http.Handle("/", handler)

	log.Printf("Starting server at http://localhost:%s", port)
	log.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed: ", err)
	}
}

package server

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Serve(port, outputDir, baseURL string) {
	// Parse the BaseURL to find the path prefix
	// Example: "https://example.com/docs" -> pathPrefix is "/docs"
	// Example: "http://localhost:8080" -> pathPrefix is ""
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Invalid BaseURL: %v", err)
	}
	pathPrefix := u.Path
	// Ensure prefix starts with / and doesn't end with / (unless it's root)
	if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
		pathPrefix = "/" + pathPrefix
	}
	pathPrefix = strings.TrimSuffix(pathPrefix, "/")

	// Setup the File Server
	fileServer := http.FileServer(http.Dir(outputDir))

	// Create the Handler
	// This logic runs AFTER the prefix has been stripped (if we use StripPrefix below)
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.URL.Path here is relative to the outputDir (prefix already stripped)
		path := r.URL.Path

		// Redirect trailing slashes to slash-less path (except root)
		if path != "/" && strings.HasSuffix(path, "/") {
			newPath := strings.TrimSuffix(path, "/")
			// We must add the prefix back for the browser redirect
			fullRedirectPath := pathPrefix + newPath
			http.Redirect(w, r, fullRedirectPath, http.StatusMovedPermanently)
			return
		}

		// Pretty URL support:
		// If requesting "/my-note", check if "/my-note.html" exists
		if filepath.Ext(path) == "" {
			// A. Check for exact HTML file (clean URLs)
			htmlPath := filepath.Join(outputDir, path+".html")
			if _, err := os.Stat(htmlPath); err == nil {
				http.ServeFile(w, r, htmlPath)
				return
			}

			// B. Check for directory index (e.g., /folder/ -> /folder/index.html)
			localPath := filepath.Join(outputDir, path)
			if info, err := os.Stat(localPath); err == nil && info.IsDir() {
				indexPath := filepath.Join(localPath, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					http.ServeFile(w, r, indexPath)
					return
				}
			}
		}

		fileServer.ServeHTTP(w, r)
	})

	// 4. Mount the handler
	// If there is a path prefix (e.g. /kiln), we strip it so the file server works,
	// but we serve it at that specific path.
	if pathPrefix != "" {
		// Handle the subpath (e.g. /kiln/)
		http.Handle(pathPrefix+"/", http.StripPrefix(pathPrefix, baseHandler))

		// Redirect root "/" to the subpath "/kiln/" for convenience
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Only redirect if it's exactly root, otherwise let it 404
			if r.URL.Path == "/" {
				http.Redirect(w, r, pathPrefix+"/", http.StatusFound)
			} else {
				http.NotFound(w, r)
			}
		})
		log.Printf("Serving at http://localhost:%s%s/", port, pathPrefix)
	} else {
		// Standard root serving
		http.Handle("/", baseHandler)
		log.Printf("Serving at http://localhost:%s", port)
	}

	log.Println("Press Ctrl+C to stop")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed: ", err)
	}
}

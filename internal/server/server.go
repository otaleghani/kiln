package server

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Serve starts a simple static file server on the specified port.
// It includes logic to handle "Clean URLs" (extensionless linking) and directory indices,
// mimicking the behavior of production static hosting providers.
func Serve(port, outputDir, baseURL string) {
	// 1. Determine Path Prefix
	// If the user's BaseURL includes a path (e.g., "https://example.com/docs"),
	// we need to serve the site under that prefix ("/docs") locally to match production.
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Invalid BaseURL: %v", err)
	}
	pathPrefix := u.Path

	// Normalize prefix: ensure it starts with "/" and doesn't end with one.
	if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
		pathPrefix = "/" + pathPrefix
	}
	pathPrefix = strings.TrimSuffix(pathPrefix, "/")

	// 2. Setup the Standard File Server
	// This will serve raw files from the output directory.
	fileServer := http.FileServer(http.Dir(outputDir))

	// 3. Create Custom Request Handler
	// This wrapper logic runs AFTER any prefix has been stripped.
	// It handles clean URLs, trailing slashes, and fallback lookups.
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.URL.Path here is relative to the outputDir (prefix already stripped by http.StripPrefix)
		path := r.URL.Path

		// A. Trailing Slash Canonicalization
		// Redirect paths ending in "/" to the non-slash version (except root).
		// Example: /about/ -> /about
		if path != "/" && strings.HasSuffix(path, "/") {
			newPath := strings.TrimSuffix(path, "/")
			// Reconstruct full path for the browser redirect (must include the prefix if one exists)
			fullRedirectPath := pathPrefix + newPath
			http.Redirect(w, r, fullRedirectPath, http.StatusMovedPermanently)
			return
		}

		// B. "Pretty URL" Support
		// If the request has no extension (e.g. "/my-note"), we try to find the actual file on disk.
		if filepath.Ext(path) == "" {
			// Case 1: Check for an HTML file with the same name.
			// Request: /my-note -> serves: /my-note.html
			htmlPath := filepath.Join(outputDir, path+".html")
			if _, err := os.Stat(htmlPath); err == nil {
				http.ServeFile(w, r, htmlPath)
				return
			}

			// Case 2: Check for a directory with an index.html.
			// Request: /folder -> serves: /folder/index.html
			localPath := filepath.Join(outputDir, path)
			if info, err := os.Stat(localPath); err == nil && info.IsDir() {
				indexPath := filepath.Join(localPath, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					http.ServeFile(w, r, indexPath)
					return
				}
			}
		}

		// Fallback: Use the standard file server (handles assets, existing files, 404s).
		fileServer.ServeHTTP(w, r)
	})

	// 4. Mount the Handler
	// If a path prefix is configured, we must strip it from the request URL
	// so that the file server sees the path relative to 'outputDir'.
	if pathPrefix != "" {
		// Handle the subpath (e.g. requests to /kiln/...)
		http.Handle(pathPrefix+"/", http.StripPrefix(pathPrefix, baseHandler))

		// Convenience: Redirect root "/" to the subpath "/kiln/"
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, pathPrefix+"/", http.StatusFound)
			} else {
				http.NotFound(w, r)
			}
		})
		log.Printf("Serving at http://localhost:%s%s/", port, pathPrefix)
	} else {
		// Standard root serving (no prefix)
		http.Handle("/", baseHandler)
		log.Printf("Serving at http://localhost:%s", port)
	}

	log.Println("Press Ctrl+C to stop")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed: ", err)
	}
}

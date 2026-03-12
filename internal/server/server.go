// Local development server with clean URL support and base path handling. @feature:dev-server
package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Serve starts a simple static file server on the specified port.
// It includes logic to handle "Clean URLs" (extensionless linking) and directory indices,
// mimicking the behavior of production static hosting providers.
func Serve(ctx context.Context, port, outputDir, baseURL string, log *slog.Logger) {
	// Determine path prefix
	// If the user's BaseURL includes a path (e.g., "https://example.com/docs"),
	// we need to serve the site under that prefix ("/docs") locally to match production.
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Error("Couldn't parse baseURL", "error", err)
		os.Exit(1)
	}
	pathPrefix := u.Path

	// Normalize prefix: ensure it starts with "/" and doesn't end with one.
	if pathPrefix != "" && !strings.HasPrefix(pathPrefix, "/") {
		pathPrefix = "/" + pathPrefix
	}
	pathPrefix = strings.TrimSuffix(pathPrefix, "/")

	// Setup the standard file server
	fileServer := http.FileServer(http.Dir(outputDir))

	notFoundPage := filepath.Join(outputDir, "404.html")

	// Create custom request handler
	// It handles clean URLs, trailing slashes, and fallback lookups
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.URL.Path here is relative to the outputDir (prefix already stripped by http.StripPrefix)
		path := r.URL.Path

		// Trailing slash canonicalization
		// Example: /about/ -> /about
		if path != "/" && strings.HasSuffix(path, "/") {
			newPath := strings.TrimSuffix(path, "/")
			// Reconstruct full path for the browser redirect (must include the prefix if one exists)
			fullRedirectPath := pathPrefix + newPath
			http.Redirect(w, r, fullRedirectPath, http.StatusMovedPermanently)
			return
		}

		// Pretty URL support
		// If the request has no extension (e.g. "/my-note"), we try to find the actual file on disk.
		if filepath.Ext(path) == "" {
			// Check for an HTML file with the same name
			// Request: /my-note -> serves: /my-note.html
			htmlPath := filepath.Join(outputDir, path+".html")
			if _, err := os.Stat(htmlPath); err == nil {
				http.ServeFile(w, r, htmlPath)
				return
			}

			// Check for a directory with an index.html (flat-urls)
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

		// Use the standard file server, but intercept 404 responses to serve
		// the custom 404.html page instead of Go's default plain text 404.
		recorder := &notFoundRecorder{ResponseWriter: w}
		fileServer.ServeHTTP(recorder, r)
		if recorder.status == http.StatusNotFound {
			serve404(w, notFoundPage)
		}
	})

	// Mount the handler
	// If a path prefix is configured, we must strip it from the request URL
	// so that the file server sees the path relative to 'outputDir'
	if pathPrefix != "" {
		// Handle the subpath (e.g. requests to /kiln/...)
		http.Handle(pathPrefix+"/", http.StripPrefix(pathPrefix, baseHandler))

		// Redirect root "/" to the subpath "/kiln/"
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, pathPrefix+"/", http.StatusFound)
			} else {
				serve404(w, notFoundPage)
			}
		})
		log.Info("Serving...", "port", port, "path", pathPrefix)
	} else {
		// Standard root serving (no prefix)
		http.Handle("/", baseHandler)
		log.Info("Serving...", "port", port)
	}

	log.Info("Press Ctrl+C to stop")

	srv := &http.Server{Addr: ":" + port}
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("Server failed", "error", err)
		os.Exit(1)
	}
	log.Info("Server stopped")
}

// serve404 writes the custom 404.html page if it exists, otherwise falls back
// to the standard Go 404 response.
func serve404(w http.ResponseWriter, notFoundPage string) {
	content, err := os.ReadFile(notFoundPage)
	if err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write(content)
}

// notFoundRecorder wraps an http.ResponseWriter to intercept 404 responses from
// the standard file server. When a 404 is detected, it suppresses the default
// response body so the caller can serve a custom 404 page instead.
type notFoundRecorder struct {
	http.ResponseWriter
	status int
}

func (r *notFoundRecorder) WriteHeader(code int) {
	r.status = code
	if code == http.StatusNotFound {
		return
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *notFoundRecorder) Write(b []byte) (int, error) {
	if r.status == http.StatusNotFound {
		return len(b), nil
	}
	return r.ResponseWriter.Write(b)
}

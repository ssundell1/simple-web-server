// customfileserver.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"simple-web-server/utils"
)

// CustomFileServer is a custom handler for serving static files
type CustomFileServer struct {
	root   http.Dir
	logger utils.Logger
}

// NewCustomFileServer creates a new CustomFileServer
func NewCustomFileServer(root string, logger utils.Logger) *CustomFileServer {
	// Check if the directory exists
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		logger.Warning("file server directory does not exist")
		logger.Info("creating directory for file server")
		err := os.MkdirAll(root, 0755) // 0755 is the default permissions
		if err != nil {
			logger.Error("failed to create file server directory")
		}
	}
	return &CustomFileServer{root: http.Dir(root), logger: logger}
}

// ServeHTTP handles the HTTP request
func (fs *CustomFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		response := "Method not allowed"
		http.Error(w, response, http.StatusMethodNotAllowed)
		fs.logger.Warning(fmt.Sprintf("%s BLOCKED %s %s", r.Method, r.RemoteAddr, r.URL))
		return
	}
	upath := r.URL.Path
	// Prevent directory traversal
	if upath == "" || upath[0] != '/' {
		upath = "/" + upath
	}
	upath = filepath.Clean(upath)

	// Serve directory listing for /files
	if upath == "/files" {
		fs.logger.Info(fmt.Sprintf("%s %s %s - Return file list", r.Method, r.RemoteAddr, r.URL))
		// Read the directory contents
		_, err := os.ReadDir("./" + upath)
		if err != nil {
			return
		}

		fs.serveDirectory(w)
		return
	}

	// Serve file contents for /files/example-file.txt
	fs.logger.Info(fmt.Sprintf("%s %s %s - Return file", r.Method, r.RemoteAddr, r.URL))
	fs.serveFile(w, r, upath)
}

// serveDirectory lists the contents of the directory
func (fs *CustomFileServer) serveDirectory(w http.ResponseWriter) {
	dir, err := fs.root.Open("/")
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		http.Error(w, "Unable to read directory contents", http.StatusInternalServerError)
		return
	}

	fileList := make([]string, 0, len(files))
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fileList); err != nil {
		http.Error(w, "Unable to encode directory contents", http.StatusInternalServerError)
	}
}

// serveFile serves the contents of a specific file
func (fs *CustomFileServer) serveFile(w http.ResponseWriter, r *http.Request, upath string) {
	// Trim the leading /files part to get the relative path
	relPath := filepath.Clean(upath[len("/files"):])

	f, err := fs.root.Open(relPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if fi.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}

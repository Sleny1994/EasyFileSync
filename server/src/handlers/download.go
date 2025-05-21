package handlers

import (
	"net/http"
	"path/filepath"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing file name", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join("uploads", name)
	http.ServeFile(w, r, filePath)
}

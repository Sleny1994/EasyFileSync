package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	relPath := r.URL.Query().Get("filename")
	if relPath == "" {
		http.Error(w, "Missing filename parameter", http.StatusBadRequest)
		return
	}
	savePath := filepath.Join("uploads", relPath)
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		http.Error(w, "Failed to create directories: "+err.Error(), http.StatusInternalServerError)
		return
	}
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Failed to write file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File %s uploaded successfully.", relPath)
}

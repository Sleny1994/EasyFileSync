package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

type FileHash struct {
	Path  string `json:"path"`
	Hash  string `json:"hash"`
	MTime int64  `json:"mtime"`
}

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	hashes, err := GetAllFileHashes("uploads")
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "Failed to read uploads", http.StatusInternalServerError)
		return
	}
	files := []FileHash{}
	for path, hash := range hashes {
		files = append(files, FileHash{Path: path, Hash: hash.Hash, MTime: hash.MTime})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

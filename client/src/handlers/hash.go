package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// 计算单个文件的MD5
func FileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 递归获取目录下所有文件的相对路径和哈希
func GetLocalFileHashes(root string) (map[string]FileHash, error) {
	hashes := make(map[string]FileHash)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, _ := filepath.Rel(root, path)
			hash, err := FileMD5(path)
			if err != nil {
				return err
			}
			hashes[rel] = FileHash{
				Path:  rel,
				Hash:  hash,
				MTime: info.ModTime().Unix(),
			}
		}
		return nil
	})
	return hashes, err
}

// 获取服务器端文件列表
func GetServerFileHashes(serverURL string) (map[string]FileHash, error) {
	resp, err := http.Get(serverURL + "/files")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var files []FileHash
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}
	result := make(map[string]FileHash)
	for _, f := range files {
		result[f.Path] = f
	}
	return result, nil
}

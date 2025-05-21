package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

// 计算指定目录下所有文件的相对路径和MD5哈希
func GetAllFileHashes(root string) (map[string]FileHash, error) {
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

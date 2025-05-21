package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileHash struct {
	Path  string `json:"path"`
	Hash  string `json:"hash"`
	MTime int64  `json:"mtime"`
}

type SyncCheckRequest struct {
	Files []FileHash `json:"files"`
}

type SyncCheckResponse struct {
	NeedSync []string `json:"need_sync"`
}

// 上传单个文件
func UploadFile(serverURL, localDir, relPath string) error {
	filePath := filepath.Join(localDir, relPath)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(relPath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", serverURL+"/upload?filename="+relPath, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed: %s", resp.Status)
	}
	return nil
}

// 下载单个文件
func DownloadFile(serverURL, localDir, relPath string) error {
	resp, err := http.Get(serverURL + "/download?name=" + relPath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	localPath := filepath.Join(localDir, relPath)
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

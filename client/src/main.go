package main

import (
	"fmt"

	"EasyFileSync/client/src/handlers"
)

func main() {
	serverURL := "http://localhost:9000"
	localDir := "./sync" // 你要同步的本地目录

	// 1. 获取本地文件元信息（含哈希和mtime）
	localMetas, err := handlers.GetLocalFileHashes(localDir)
	if err != nil {
		fmt.Println("本地哈希计算失败:", err)
		return
	}

	// 2. 获取服务端文件元信息（含哈希和mtime）
	serverMetas, err := handlers.GetServerFileHashes(serverURL)
	if err != nil {
		fmt.Println("获取服务端文件列表失败:", err)
		return
	}

	// 3. 比较并同步
	// 3.1 本地有的文件
	for path, local := range localMetas {
		server, exists := serverMetas[path]
		switch {
		case !exists:
			// 服务端没有，上传
			fmt.Printf("上传: %s ... ", path)
			if err := handlers.UploadFile(serverURL, localDir, path); err != nil {
				fmt.Println("失败:", err)
			} else {
				fmt.Println("成功")
			}
		case local.Hash != server.Hash:
			if local.MTime > server.MTime {
				// 本地新，上传
				fmt.Printf("上传: %s ... ", path)
				if err := handlers.UploadFile(serverURL, localDir, path); err != nil {
					fmt.Println("失败:", err)
				} else {
					fmt.Println("成功")
				}
			} else if local.MTime < server.MTime {
				// 服务端新，下载
				fmt.Printf("从服务端下载: %s ... ", path)
				if err := handlers.DownloadFile(serverURL, localDir, path); err != nil {
					fmt.Println("失败:", err)
				} else {
					fmt.Println("成功")
				}
			} else {
				// mtime 相等但内容不同，冲突
				fmt.Printf("冲突文件: %s，请手动处理\n", path)
			}
		}
	}
	// 3.2 服务端有但本地没有的文件
	for path := range serverMetas {
		if _, exists := localMetas[path]; !exists {
			fmt.Printf("从服务端下载: %s ... ", path)
			if err := handlers.DownloadFile(serverURL, localDir, path); err != nil {
				fmt.Println("失败:", err)
			} else {
				fmt.Println("成功")
			}
		}
	}
}

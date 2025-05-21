package main

import (
	"fmt"
	"log"
	"net/http"

	"EasyFileSync/server/src/handlers"
)

func main() {
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/files", handlers.FilesHandler)
	http.HandleFunc("/download", handlers.DownloadHandler)

	go func() {
		fmt.Println("HTTP Web Service is listening on localhost:9000")
		log.Fatal(http.ListenAndServe(":9000", nil))
	}()

	// https 检查证书文件是否存在
	// certFile := filepath.Join(".", "server.crt")
	// keyFile := filepath.Join(".", "server.key")
	// if _, err := os.Stat(certFile); err == nil {
	// 	if _, err := os.Stat(keyFile); err == nil {
	// 		fmt.Println("HTTPS 服务启动于 :9443")
	// 		log.Fatal(http.ListenAndServeTLS(":9443", certFile, keyFile, nil))
	// 	}
	// }

	select {} // 阻塞主线程
}

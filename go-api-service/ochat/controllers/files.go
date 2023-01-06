package controllers

import (
	"io"
	"net/http"
	"ochat/comm/funcs"
	"os"
)

// 显示文件
func ImageShow(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	path := r.FormValue("path")
	if filename == "" || path == "" {
		http.NotFound(w, r)
		return
	}

	filePath := funcs.GetUploadFilePath(path, filename)

	fileBytes, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer fileBytes.Close()

	io.Copy(w, fileBytes)
}

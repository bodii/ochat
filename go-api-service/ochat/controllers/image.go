package controllers

import (
	"io"
	"net/http"
	"ochat/bootstrap"
	"ochat/comm"
	"os"
	"strings"
)

// 显示图片
func ImgShow(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	fileNameIndex := strings.LastIndexByte(path, '/')
	file := path[fileNameIndex+1:]
	suffixIndex := strings.LastIndexByte(file, '.')
	if suffixIndex < 1 {
		comm.Res(w, 404, "file not exists!", nil)
		return
	}

	filename, suffix := file[:suffixIndex], file[suffixIndex:]
	avatarConf := bootstrap.SystemConf.Avatar
	file_byte, err := os.Open(
		comm.GetProjectDIR() + avatarConf.FileDir + filename + suffix)
	if err != nil {
		comm.Res(w, 1001, err.Error(), nil)
	}
	defer file_byte.Close()

	io.Copy(w, file_byte)
}

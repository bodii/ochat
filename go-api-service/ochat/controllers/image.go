package controllers

import (
	"io"
	"net/http"
	"ochat/bootstrap"
	"ochat/comm"
	"os"
	"strings"
)

// 显示头像图片
func ShowAvatar(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	defer file_byte.Close()

	io.Copy(w, file_byte)
}

// 上传头像图片
func UpPicture(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("img")
	if err != nil {
		comm.Res(w, 1001, err.Error(), nil)
		return
	}
	defer file.Close()

	filename := fileHeader.Filename
	suffix := strings.SplitAfter(fileHeader.Filename, ".")
	name := comm.GetRandFileName(suffix[1])
	avatarConf := bootstrap.SystemConf.Avatar
	filepath := bootstrap.PROJECT_DIR + avatarConf.UploadDir + name

	savePath, err := os.Create(filepath)
	if err != nil {
		comm.Res(w, 1001, err.Error(), nil)
		return
	}

	io.Copy(savePath, file)

	picUrl := bootstrap.HTTP_Avatar_URI + filename
	result := map[string]any{"pic_url": picUrl}
	comm.Res(w, 100, "update success!", result)
}

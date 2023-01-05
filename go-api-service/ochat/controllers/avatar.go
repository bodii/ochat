package controllers

import (
	"io"
	"net/http"
	"ochat/bootstrap"
	"ochat/comm"
	"ochat/comm/funcs"
	"os"
	"path/filepath"
	"strings"
)

// 显示头像图片
func AvatarShow(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	fileNameIndex := strings.LastIndexByte(path, '/')
	file := path[fileNameIndex+1:]
	suffixIndex := strings.LastIndexByte(file, '.')
	if suffixIndex < 1 {
		comm.ResFailure(w, 404, "file not exists!")
		return
	}

	filename, suffix := file[:suffixIndex], file[suffixIndex:]
	avatarConf := bootstrap.SystemConf.Avatar
	file_byte, err := os.Open(
		funcs.GetProjectDIR() + avatarConf.FileDir + filename + suffix)
	if err != nil {
		comm.ResFailure(w, 1001, err.Error())
		return
	}
	defer file_byte.Close()

	io.Copy(w, file_byte)
}

// 上传头像图片
func AvatarUpload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("img")
	if err != nil {
		comm.ResFailure(w, 1001, err.Error())
		return
	}
	defer file.Close()

	filename := fileHeader.Filename
	name := funcs.RandFileName(filepath.Ext(filename))
	avatarConf := bootstrap.SystemConf.Avatar
	filepath := bootstrap.PROJECT_DIR + avatarConf.UploadDir + name

	savePath, err := os.Create(filepath)
	if err != nil {
		comm.ResFailure(w, 1001, err.Error())
		return
	}

	io.Copy(savePath, file)

	picUrl := bootstrap.HTTP_Avatar_URI + filename
	result := map[string]any{"pic_url": picUrl}
	comm.ResSuccess(w, result)
}

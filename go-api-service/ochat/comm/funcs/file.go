package funcs

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"ochat/bootstrap"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/skip2/go-qrcode"
)

func GetProjectDIR() string {
	dir, err := os.Getwd()

	if err != nil {
		return os.Args[0]
	}

	return dir
}

// RandFileName func
//
// param:
//   - suffix the affix of the file name to be generated.
//
// returns a file name randomly generated by suffix.
func RandFileName(suffix string) string {
	return fmt.Sprintf("%d%04d%s",
		time.Now().Unix(),
		GetUnixNanoRandSeed().Int31(),
		suffix)
}

// create QrCode func
//
// param:
//   - url [string]: QR Code to url addr.
//   - pathTag [string]: QR Code save path the tag.
//
// return:
//   - filename [string]: QR Code filename
//   - err [error]: error info
func QrCode(url string, pathTag string) (filename string, err error) {
	filename = RandFileName(".png")
	filePath := path.Join(GetProjectDIR(), "/files/upload/", pathTag, filename)
	err = qrcode.WriteFile(url, qrcode.Medium, 256, filePath)
	if err != nil {
		log.Printf(" %s create %s failure\n", url, pathTag)
		return filename, err
	}

	return filename, nil
}

// copy file func
//
// param:
//   - dst [string]: the target path to save the file
//   - src [string]: source storage path of the file
//
// return:
//   - err [error]: error info
func CopyFile(dst string, src string) error {
	// open src file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	// create dst file
	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	// last close files
	defer srcFile.Close()
	defer dstFile.Close()

	// copy
	_, err = io.Copy(dstFile, srcFile)

	// return
	return err
}

// get the url of the image generated
//
// param:
//   - pathTag [string]: tag of the save path
//   - filename [string]: filename of the save file
//
// return:
//   - url [string]: the url address is generated successfully
func GetImgUrl(pathTag, filename string) (url string) {
	if pathTag == "" || filename == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s?path=%s&filename=%s",
		bootstrap.HTTP_HOST, "files/image", pathTag, filename)
}

// get a default profile picture based on gender
//
// param:
//   - sex [int]: sex value
//
// return:
//   - filename [string]: filename of the save file
func DefaultAvatar(sex int) (filename string) {
	staticAvatarPath := path.Join(GetProjectDIR(), "/files/static/avatar/default/")
	defaultAvatar := ""
	switch sex {
	case 1:
		defaultAvatar = "avatar_boy_kid_person_icon.png"
	case 2:
		defaultAvatar = "child_girl_kid_person_icon.png"
	default:
		defaultAvatar = "avatar_boy_male_user_young_icon.png"
	}
	staticAvatarPath = path.Join(staticAvatarPath, defaultAvatar)

	newFilename := RandFileName(".png")
	newAvatarPath := GetUploadFilePath("avatar", newFilename)
	// copy file
	err := CopyFile(newAvatarPath, staticAvatarPath)
	if err != nil {
		return ""
	}

	return newFilename
}

// get upload file path
//
// param:
//   - pathTag [string]: tag of the save path
//   - filename [string]: filename of the save file
//
// return:
//   - [string]: path to save the file
func GetUploadFilePath(pathTag, fielname string) string {
	return path.Join(GetProjectDIR(), "/files/upload/", pathTag, fielname)
}

// get a default profile picture based on gender
//
// param:
//   - sex [int]: sex value
//
// return:
//   - filename [string]: filename of the save file
func UploadFile(r *http.Request, upName, pathTag string) (filename, oldFilename string, err error) {
	file, fileHeader, err := r.FormFile(upName)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	oldFilename = fileHeader.Filename
	filename = RandFileName(filepath.Ext(oldFilename))
	filepath := GetUploadFilePath(pathTag, filename)

	savePath, err := os.Create(filepath)
	if err != nil {
		return filename, oldFilename, err
	}
	defer savePath.Close()

	io.Copy(savePath, file)

	return
}

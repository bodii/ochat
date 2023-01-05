package funcs

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// MD5 result a data string to md5 value
func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

// MD5Enhance func
// return a enhance n level the md5 value
func MD5Enhance(data string, n int) string {
	h := md5.New()
	h.Write([]byte(data))
	str := h.Sum(nil)
	for i := 1; i < n; i++ {
		h.Reset()
		h.Write(str)
		str = h.Sum(nil)
	}

	return hex.EncodeToString(str)
}

// md5转大写
func MD5ToUpper(data string) string {
	return strings.ToUpper(MD5(data))
}

// 获取纳秒生成md5字符
func GetNano2md5String() string {
	r := GetUnixNanoRandSeed()
	number := r.Int63n(time.Now().UnixNano())
	return MD5(strconv.FormatInt(number, 10))
}

// 生成一个加盐的密码
func GeneratePasswd(pwd, salt string) string {
	return MD5(pwd + salt)
}

// 获取一个token, 按小时生效
func GenerateToken(data string) string {
	nowStr, _ := time.Parse("2006010215", time.Now().String())

	return MD5ToUpper(data + nowStr.String())
}

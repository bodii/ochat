package comm

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

var lowletters = "abcdefghijklmnopqrstuvwxyz"
var upLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var numbers = "0123456789"

const (
	// 小写字母[a-z], length: 26
	RandStrlevel1 int = iota + 1
	// 大写字母[A-Z], length: 26
	RandStrlevel2
	// 数字字符串[01...9], length: 10
	RandStrlevel3
	// 小写字母[a-z] + 大写字母[A-Z], length: 52
	RandStrlevel4
	// 小写字母[a-z] + 大写字母[A-Z] + 数字[01...9], length: 62
	RandStrlevel5
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

func MD5ToUpper(data string) string {
	return strings.ToUpper(MD5(data))
}

func GeneratePasswd(pwd, salt string) string {
	return MD5(pwd + salt)
}

func VaildataPasswd(newPwd, salt, oldPwdAndSalt string) bool {
	return GeneratePasswd(newPwd, salt) == oldPwdAndSalt
}

// func RandStr
// params: [length int] rand string length
// [level int] 1: 小写字母；2：大写字母；3：数字；4：小写字母+大写字母；5：小写字母+大写字母+数字
// returns: a random string of length and level
func RandStr(length int, level int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	letters := randstr_level2string(level)
	lettersLen := len(letters)

	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(lettersLen)]
	}

	return string(b)
}

func randstr_level2string(level int) string {
	switch level {
	case RandStrlevel1:
		return lowletters
	case RandStrlevel2:
		return upLetters
	case RandStrlevel3:
		return numbers
	case RandStrlevel4:
		return lowletters + upLetters
	case RandStrlevel5:
		return lowletters + upLetters + numbers
	default:
		return lowletters
	}
}

func GenerateToken(data string) string {
	nowStr, _ := time.Parse("2006010215", time.Now().String())

	return MD5ToUpper(data + nowStr.String())
}

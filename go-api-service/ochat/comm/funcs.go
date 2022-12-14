package comm

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	// 小写字母
	low_letters = "abcdefghijklmnopqrstuvwxyz"
	// 大写字母
	up_letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 数字字符串
	numbers = "0123456789"
	// 特殊字符
	special_characters = "-+=<>@.;^*()~{}[]|:/"
)

const (
	// 小写字母[a-z], length: 26
	Rand_Str_Level_1 int = iota + 1
	// 大写字母[A-Z], length: 26
	Rand_Str_Level_2
	// 数字字符串[01...9], length: 10
	Rand_Str_Level_3
	// 小写字母[a-z] + 大写字母[A-Z], length: 52
	Rand_Str_Level_4
	// 小写字母[a-z] + 大写字母[A-Z] + 数字[01...9], length: 62
	Rand_Str_Level_5
	// 小写字母[a-z] + 大写字母[A-Z] + 数字[01...9] + 特殊字符, length: 82
	Rand_Str_Level_6
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

// 生成一个加盐的密码
func GeneratePasswd(pwd, salt string) string {
	return MD5(pwd + salt)
}

// 验证密码是否正确
func VaildataPasswd(newPwd, salt, oldPwdAndSalt string) bool {
	return GeneratePasswd(newPwd, salt) == oldPwdAndSalt
}

// 获取纳秒的随机种子
func GetUnixNanoRandSeed() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 获取纳秒生成md5字符
func GetNano2md5String() string {
	r := GetUnixNanoRandSeed()
	number := r.Int63n(time.Now().UnixNano())
	return MD5(strconv.FormatInt(number, 10))
}

// func RandStr
// params: [length int] rand string length
// [level int] 1: 小写字母；2：大写字母；3：数字；4：小写字母+大写字母；5：小写字母+大写字母+数字
// returns: a random string of length and level
func RandStr(length int, level int) string {
	r := GetUnixNanoRandSeed()

	letters := getStrLevel2string(level)
	lettersLen := len(letters)

	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(lettersLen)]
	}

	return string(b)
}

// 获取对应等级的字符串值
func getStrLevel2string(level int) string {
	switch level {
	case Rand_Str_Level_1:
		return low_letters
	case Rand_Str_Level_2:
		return up_letters
	case Rand_Str_Level_3:
		return numbers
	case Rand_Str_Level_4:
		return low_letters + up_letters
	case Rand_Str_Level_5:
		return low_letters + up_letters + numbers
	case Rand_Str_Level_6:
		return low_letters + up_letters + numbers + special_characters
	default:
		return low_letters
	}
}

// 获取一个token, 按小时生效
func GenerateToken(data string) string {
	nowStr, _ := time.Parse("2006010215", time.Now().String())

	return MD5ToUpper(data + nowStr.String())
}

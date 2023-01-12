package funcs

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// get a string md5 value
//
// param:
//   - data [string]: a string data
//
// return:
//   - [string]: the generated md5 string
func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

// generate md5 value after many times
//
// param:
//   - data [string]: a string data
//   - n [int]: number of times to md5
//
// return:
//   - [string]: the generated md5 string
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

// generates the uppercase value of the string after md5
//
// param:
//   - data [string]: a string data
//
// return:
//   - [string]: the generated upper md5 string
func MD5ToUpper(data string) string {
	return strings.ToUpper(MD5(data))
}

// generates an md5 encrypted value for nanoseconds
//
// param:
//
// return:
//   - [string]: the generated md5 string
func GetNano2md5String() string {
	r := GetUnixNanoRandSeed()
	number := r.Int63n(time.Now().UnixNano())
	return MD5(strconv.FormatInt(number, 10))
}

// generate a salted password
//
// param:
//   - pwd [string]: password value
//   - salt [string]: the value of the salt property
//
// return:
//   - [string]: the generated md5 string
func GeneratePasswd(pwd, salt string) string {
	return MD5(pwd + salt)
}

// generates an uppercase token value for the current hour
//
// param:
//   - data [string]: a string value
//
// return:
//   - [string]: the generated upper token encryted value
func GenerateToken(data string) string {
	nowStr, _ := time.Parse("2006010215", time.Now().String())

	return MD5ToUpper(data + nowStr.String())
}

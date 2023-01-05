package funcs

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
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

// get str prefix
//
// param:
//   - str string.
//   - length int.
//   - p_type int. 1:默认;2:大写
//
// returns str 指定长度的前缀.
func StrPrefix(str string, length int, p_type int) string {
	strRune := []rune(str)
	strPrefix := string(strRune[:length])

	enStr := ""
	if IsChinese(strPrefix) {
		pyStrLists := pinyin.Pinyin(strPrefix, pinyin.NewArgs())
		enStrList := make([]string, len(pyStrLists))
		for i, en := range pyStrLists {
			enStrList[i] = en[0]
		}

		enStr = strings.Join(enStrList, "")[:length]
	}

	if IsEnglish(strPrefix) && p_type == 2 {
		enStr = strings.ToUpper(strPrefix)
	}

	return enStr
}

// ISChinese func
//
// param:
//   - str string.
//
// returns is str chinese true | false
func IsChinese(str string) bool {
	return regexp.MustCompile("^[\u4e00-\u9fa5]$").MatchString(str)
}

// IsEnglish func
//
// param:
//   - str string.
//
// returns is str english true | false
func IsEnglish(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(str)
}

// IsNumber func
//
// param: str string
//
// returns: this is a number: true|false
func IsNumber(str string) bool {
	return regexp.MustCompile(`^\d+$`).Match([]byte(str))
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

// func RandStr
//
//	 params:
//
//	[ length int ]
//
//		rand string length
//
//	[ level int ]
//	  1：小写字母；
//	  2：大写字母；
//	  3：数字；
//	  4：小写字母+大写字母；
//	  5：小写字母+大写字母+数字;
//	  6：小写字母+大写字母+数字+特殊字符
//
// <= returns :
//
//	a random string of length and level
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

package funcs

import "regexp"

// 验证密码是否正确
func VaildataPasswd(newPwd, salt, oldPwdAndSalt string) bool {
	return GeneratePasswd(newPwd, salt) == oldPwdAndSalt
}

// IsMobile func
//
// param: str string
//
// returns this is a mobile: true|false
func IsMobile(str string) bool {
	return regexp.MustCompile(`^(1[3|4|5|7|8|9]\d{9})$`).
		Match([]byte(str))
}

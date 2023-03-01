package funcs

import "regexp"

// verify that the password is valid
//
// param:
//   - newPwd [string]: input password value
//   - salt [string]: the value of the salt property
//   - oldPwdAndSalt [string]: contains the password after salt and password
//
// return:
//   - [bool]: determine whether the new salted password is the same as the old password after encryption
func VaildataPasswd(newPwd, salt, oldPwdAndSalt string) bool {
	return GeneratePasswd(newPwd, salt) == oldPwdAndSalt
}

// verify that the string value is a cell phone number
//
// param:
//   - str [string]: a string value that may contain the phone number
//
// return:
//   - [bool] this is a mobile: true|false
func IsMobile(str string) bool {
	if str == "" {
		return false
	}

	return regexp.MustCompile(`^(1[3|4|5|7|8|9]\d{9})$`).
		Match([]byte(str))
}

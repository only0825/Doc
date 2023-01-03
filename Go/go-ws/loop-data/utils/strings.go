package utils

import (
	"math/rand"
	"time"
	"unicode"
)

func CtypePunct(s string) bool {

	if s == "" {
		return false
	}
	for _, r := range s {
		if unicode.IsPunct(r) || unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

// 判断字符串是不是数字
func CtypeDigit(s string) bool {

	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// 判断字符串是不是字母+数字
func CtypeAlnum(s string) bool {

	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// 获取source的子串,如果start小于0或者end大于source长度则返回""
// start:开始index，从0开始，包括0
// end:结束index，以end结束，但不包括end
func Substring(source string, start int, end int) string {

	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

// 字符串特殊字符转译
func Addslashes(str string) string {

	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}

	return string(tmpRune)
}

func RandomKey(length int) string {

	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

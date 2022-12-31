package validator

import (
	"github.com/shopspring/decimal"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	loc, _   = time.LoadLocation("Asia/Shanghai")
	cuprefix = map[string]bool{
		"a1": true,
		"b1": true,
		"c1": true,
	}
	caprefix = map[string]bool{
		"a": true,
		"b": true,
		"c": true,
	}
	cu = map[int]string{
		0: "a1",
		1: "b1",
		2: "c1",
	}
	ca = map[int]string{
		0: "a",
		1: "b",
		2: "c",
	}
	phoneMp = map[string]bool{
		//Viettel 086,096,097,098,032,033,034,035,036,037,038,039
		"086": true,
		"096": true,
		"097": true,
		"098": true,
		"032": true,
		"033": true,
		"034": true,
		"035": true,
		"036": true,
		"037": true,
		"038": true,
		"039": true,
		//Vina Phone 091,094,083,084,085,081,082,088
		"091": true,
		"094": true,
		"083": true,
		"084": true,
		"085": true,
		"081": true,
		"082": true,
		"088": true,
		//Mobifone 089,090,093,070,079,077,076,078
		"089": true,
		"090": true,
		"093": true,
		"070": true,
		"079": true,
		"077": true,
		"076": true,
		"078": true,
		//Vietnamobile 092, 056, 058, 052
		"092": true,
		"056": true,
		"058": true,
		"052": true,
		//Gmobie 099, 059
		"099": true,
		"059": true,
		// Itelecom 087
		"087": true,
	}
	zaloMp = map[string]bool{
		//Viettel 086,096,097,098,032,033,034,035,036,037,038,039
		"86": true,
		"96": true,
		"97": true,
		"98": true,
		"32": true,
		"33": true,
		"34": true,
		"35": true,
		"36": true,
		"37": true,
		"38": true,
		"39": true,
		//Vina Phone 091,094,083,084,085,081,082,088
		"91": true,
		"94": true,
		"83": true,
		"84": true,
		"85": true,
		"81": true,
		"82": true,
		"88": true,
		//Mobifone 089,090,093,070,079,077,076,078
		"89": true,
		"90": true,
		"93": true,
		"70": true,
		"79": true,
		"77": true,
		"76": true,
		"78": true,
		//Vietnamobile 092, 056, 058, 052
		"092": true,
		"056": true,
		"058": true,
		"052": true,
		//Gmobie 099, 059
		"99": true,
		"59": true,
		// Itelecom 087
		"87": true,
	}
)

// 判断字符是否为数字
func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

// 判断字符是否为英文字符
func isAlpha(r rune) bool {

	if r >= 'A' && r <= 'Z' {
		return true
	} else if r >= 'a' && r <= 'z' {
		return true
	}
	return false
}

func CheckStringVName(s string) bool {

	if s == "" {
		return false
	}

	nums := 0
	for _, r := range s {
		if r == ' ' {
			nums++
		}
		if (r < 'A' || r > 'Z') && r != ' ' {
			return false
		}
	}

	if nums < 1 || nums > 4 {
		return false
	}
	return true
}

func isPriv(s string) bool {

	if s == "" {
		return false
	}

	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && r != '_' {
			return false
		}
	}

	return true
}

// 检测会员名
func CheckUName(str string, min, max int) bool {

	if !CtypeAlnum(str) || //数字字母组合
		!FirstIsAlpha(str) || //必须包含字母
		!CheckStringLength(str, min, max) {
		return false
	}

	return true
}

// 检测后台账号
func CheckAName(str string, min, max int) bool {

	if !CtypeAlnum(str) || //数字字母组合
		!FirstIsAlpha(str) || //必须包含字母
		!CheckStringLength(str, min, max) {
		return false
	}

	return true
}

// 检测信用盘会员名格式
func CheckCUName(str string, min, max int) bool {

	str = strings.ToLower(str)
	if !CheckStringLength(str, min, max) {
		return false
	}
	// 前缀不正确
	_, ok := cuprefix[str[:2]]
	if !ok {
		return false
	}

	if !CtypeAlnum(str) { //数字字母组合
		return false
	}

	return true
}

// 检测添加信用盘会员名格式
func CheckACUName(str string, level, min, max int) bool {

	str = strings.ToLower(str)
	// 层级不支持
	prefix, ok := cu[level]
	if !ok {
		return false
	}

	if !CtypeAlnum(str) || //数字字母组合
		!strings.HasPrefix(str, prefix) || //必须包含层级会员前缀
		!CheckStringLength(str, min, max) {
		return false
	}

	return true
}

// 检测添加信用盘代理名格式
func CheckCAName(str string, min, max int) bool {

	str = strings.ToLower(str)
	if !CheckStringLength(str, min, max) {
		return false
	}

	// 前缀不正确
	_, ok := caprefix[str[:1]]
	if !ok {
		return false
	}

	if !CtypeAlnum(str) { //数字字母组合
		return false
	}

	return true
}

// 检测添加信用盘代理名格式
func CheckACAName(str string, level, min, max int) bool {

	str = strings.ToLower(str)
	// 层级不支持
	uPrefix, ok := cu[level]
	if !ok {
		return false
	}

	// 层级不支持
	aPrefix, ok := ca[level]
	if !ok {
		return false
	}

	if !CtypeAlnum(str) || //数字字母组合
		!strings.HasPrefix(str, aPrefix) || //必须包含层级代理前缀
		strings.HasPrefix(str, uPrefix) || //不能包含层级会员前缀
		!CheckStringLength(str, min, max) {
		return false
	}

	return true
}

// 检测会员密码
func CheckUPassword(str string, min, max int) bool {

	if !CtypeAlnum(str) || //数字字母组合
		!CheckStringLength(str, min, max) ||
		!IncludeAlpha(str) || //必须包含字母
		!IncludeDigit(str) { //必须包含数字
		return false
	}

	return true
}

// 检测后台密码（会员密码和后台密码暂时规则相同，单独函数方便以后扩展）
func CheckAPassword(str string, min, max int) bool {

	if !CtypeAlnum(str) || //数字字母组合
		!CheckStringLength(str, min, max) ||
		!IncludeAlpha(str) || //必须包含字母
		!IncludeDigit(str) { //必须包含数字
		return false
	}

	return true
}

// 匹配值是否为空
func checkStr(str string) bool {

	n := len(str)
	if n <= 0 {
		return false
	}

	return true
}

// 判断是否为bool
func checkBool(str string) bool {

	_, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return true
}

// 判断是否为float
func CheckFloat(str string) bool {

	_, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return false
	}
	return true
}

// 判断长度
func checkLength(str string, min, max int) bool {

	if min == 0 && max == 0 {
		return true
	}

	n := len(str)
	if n < min || n > max {
		return false
	}

	return true
}

// 判断字符串长度
func CheckStringLength(val string, _min, _max int) bool {

	if _min == 0 && _max == 0 {
		return true
	}

	count := utf8.RuneCountInString(val)
	if count < _min || count > _max {

		return false
	}
	return true
}

// 判断数字范围
func CheckIntScope(s string, min, max int64) bool {

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false
	}

	if val < min || val > max {
		return false
	}

	return true
}

// 判断浮点范围
func CheckFloatScope(s, min, max string) (decimal.Decimal, bool) {

	fs, err := decimal.NewFromString(s)
	if err != nil {
		return fs, false
	}

	fMin, err := decimal.NewFromString(min)
	if err != nil {
		return fs, false
	}

	fMax, err := decimal.NewFromString(max)
	if err != nil {
		return fs, false
	}

	if fs.Cmp(fMin) == -1 || fMax.Cmp(fs) == -1 {
		return fs, false
	}

	return fs, true
}

// 判断金额上下边界
func CheckAmountRange(low, up string) (int, error) {

	fLow, err := decimal.NewFromString(low)
	if err != nil {
		return 0, err
	}

	fUp, err := decimal.NewFromString(up)
	if err != nil {
		return 0, err
	}

	return fUp.Cmp(fLow), nil
}

// 判断是否全为数字
func CheckStringDigit(s string) bool {

	if s == "" {
		return false
	}
	for _, r := range s {
		if (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}

// 判断是否全为数字+逗号
func CheckStringCommaDigit(s string) bool {

	if s == "" {
		return false
	}
	for _, r := range s {
		if (r < '0' || r > '9') && r != ',' && r != '|' {
			return false
		}
	}
	return true
}

// 判断是不是中文
func CheckStringCHN(str string) bool {

	for _, r := range str {
		if !unicode.Is(unicode.Han, r) &&
			!isAlpha(r) && (r < '0' || r > '9') && r != '_' &&
			r != ' ' && r != '-' && r != '!' && r != '@' && r != ':' &&
			r != '?' && r != '+' && r != '.' && r != '/' && r != '\'' &&
			r != '(' && r != ')' && r != '·' && r != '&' {
			return false
		}
	}
	return true
}

// 判断是不是英文数字或者汉字
func CheckStringCHNAlnum(str string) bool {

	for _, r := range str {
		if !isDigit(r) && !isAlpha(r) &&
			r != ' ' && r != '-' && r != '!' && r != '_' &&
			r != '@' && r != '?' && r != '+' && r != ':' &&
			r != '.' && r != '/' && r != '(' && r != '\'' &&
			r != ')' && r != '·' && r != '&' && !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}

// 判断是否module格式
func CheckStringModule(s string) bool {

	if s == "" {
		return false
	}

	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && r != '/' {
			return false
		}
	}

	return true
}

// 判断是否全英文字母
func CheckStringAlpha(s string) bool {

	if s == "" {
		return false
	}

	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && r != ' ' {
			return false
		}
	}

	return true
}

// 判断是否全英文字母+逗号
func CheckStringCommaAlpha(s string) bool {

	if s == "" {
		return false
	}

	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < '0' || r > '9') && (r < 'a' || r > 'z') && r != ',' {
			return false
		}
	}

	return true
}

// 判断是否全为英文字母和数字组合
func CheckStringAlnum(s string) bool {

	if s == "" {
		return false
	}
	for _, r := range s {
		if !isDigit(r) && !isAlpha(r) &&
			r != ' ' && r != '-' && r != '!' && r != '_' &&
			r != '@' && r != '?' && r != '+' && r != ':' &&
			r != '.' && r != '/' && r != '(' && r != '\'' &&
			r != ')' && r != '·' && r != '&' {
			return false
		}
	}
	return true
}

// 检查日期格式"YYYY-MM-DD"
func CheckDate(str string) bool {

	_, err := time.ParseInLocation("2006-01-02", str, loc)
	if err != nil {
		return false
	}
	return true
}

// 匹配时间 "HH:ii" or "HH:ii:ss"
func checkTime(str string) bool {

	_, err := time.ParseInLocation("15:04:05", str, loc)
	if err != nil {
		return false
	}
	return true
}

// 检查日期时间格式"YYYY-MM-DD HH:ii:ss"
func CheckDateTime(str string) bool {

	_, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	if err != nil {
		return false
	}
	return true
}

func CheckMoney(money string) bool {

	// 金额小数验证
	_, err := strconv.Atoi(money)
	if err != nil {
		return false
	}
	_, err = strconv.ParseFloat(money, 64)
	if err != nil {
		return false
	}
	return true
}

// 判断字符串是不是数字
func CtypeDigit(s string) bool {

	if s == "" {
		return false
	}
	for _, r := range s {
		if !isDigit(r) {
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
		if !isDigit(r) && !isAlpha(r) {
			return false
		}
	}

	return true
}

// 判断是否包含数字
func IncludeDigit(s string) bool {

	if s == "" {
		return false
	}

	for _, r := range s {
		if isDigit(r) {
			return true
		}
	}

	return false
}

// 判断是否包含字母
func IncludeAlpha(s string) bool {

	if s == "" {
		return false
	}

	for _, r := range s {
		if isAlpha(r) {
			return true
		}
	}

	return false
}

// 判断字符串是不是字母开头
func FirstIsAlpha(s string) bool {

	if s == "" {
		return false
	}

	r := []rune(s)
	if len(r) < 2 {
		return false
	}

	if !isAlpha(r[0]) {
		return false
	}

	return true
}

// 检查url
func CheckUrl(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func zip(a1, a2 []string) []string {

	r := make([]string, 2*len(a1))
	for i, e := range a1 {
		r[i*2] = e
		r[i*2+1] = a2[i]
	}

	return r
}

func UnFilter(str string) string {

	array2 := []string{"<", ">", "&", `"`, " "}
	array1 := []string{"&lt;", "&gt;", "&amp;", "&quot;", "&nbsp;"}

	return strings.NewReplacer(zip(array1, array2)...).Replace(str)
}

func FilterInjection(str string) string {

	array1 := []string{"<", ">", "&", `"`, " "}
	array2 := []string{"&lt;", "&gt;", "&amp;", "&quot;", "&nbsp;"}

	return strings.NewReplacer(zip(array1, array2)...).Replace(str)
}

func IsVietnamesePhone(phone string) bool {

	if !CtypeDigit(phone) {
		return false
	}

	if len(phone) != 10 {
		return false
	}

	if !strings.HasPrefix(phone, "0") {
		return false
	}

	prefix := phone[:3]
	if _, ok := phoneMp[prefix]; ok {
		return true
	}

	return false
}

func IsVietnameseZalo(zalo string) bool {

	if !CtypeDigit(zalo) {
		return false
	}

	if len(zalo) != 9 {
		return false
	}

	if strings.HasPrefix(zalo, "0") {
		return false
	}

	prefix := zalo[:2]
	if _, ok := zaloMp[prefix]; ok {
		return true
	}

	return false
}

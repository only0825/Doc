package helper

import (
	"encoding/hex"
	"fmt"
	"unicode"
	"github.com/minio/md5-simd"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
	"lukechampine.com/frand"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/*
const (
	TYpeNone = iota
	TypePhone
	TypeBankCardNumber
	TypeVirtualCurrencyAddress
	TypeRealName
	TypeEmail
)
*/

type Response struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

func Print(ctx *fasthttp.RequestCtx, state bool, data interface{}) {

	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")

	res := Response{
		Status: state,
		Data:   data,
	}

	bytes, err := JsonMarshal(res)
	if err != nil {
		ctx.SetBody([]byte(err.Error()))
		return
	}

	ctx.SetBody(bytes)
}

func PrintJson(ctx *fasthttp.RequestCtx, state bool, data string) {

	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")

	builder := strings.Builder{}

	builder.WriteString(`{"status":`)
	builder.WriteString(strconv.FormatBool(state))
	builder.WriteString(`,"data":`)
	builder.WriteString(data)
	builder.WriteString("}")

	ctx.SetBody([]byte(builder.String()))
}

func GenId() string {

	var min uint64 = 0
	var max uint64 = 9

	return fmt.Sprintf("%d%d", Cputicks(), frand.Uint64n(max-min)+min)
}

func GenLongId() string {

	var min uint64 = 100000
	var max uint64 = 999999

	id := fmt.Sprintf("%d%d", Cputicks(), frand.Uint64n(max-min)+min)
	return id[0:18]
}



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

//判断字符串是不是数字
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

//判断字符串是不是字母+数字
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

//获取source的子串,如果start小于0或者end大于source长度则返回""
//start:开始index，从0开始，包括0
//end:结束index，以end结束，但不包括end
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

//字符串特殊字符转译
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

func GetMD5Hash(text string) string {

	server := md5simd.NewServer()
	md5Hash := server.NewHash()

	md5Hash.Write([]byte(text))

	digest := md5Hash.Sum([]byte{})
	encrypted := hex.EncodeToString(digest)

	server.Close()
	md5Hash.Close()

	return encrypted
}

func TrimStr(val decimal.Decimal) string {

	s := "0.000"
	sDigit := strings.Split(val.String(), ".")
	if len(sDigit) != 2 {
		if len(sDigit) == 1 && CtypeDigit(sDigit[0]) {
			return sDigit[0] + ".000"
		}
		return s
	}

	// 浮点位数校验
	if len(sDigit[1]) <= 3 {
		return val.String()
	}

	return sDigit[0] + "." + sDigit[1][:3]
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

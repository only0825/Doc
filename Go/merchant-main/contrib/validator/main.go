package validator

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/modern-go/reflect2"
	"github.com/valyala/fasthttp"
)

func Bind(ctx *fasthttp.RequestCtx, objs interface{}) error {

	rt := reflect2.TypeOf(objs)
	rtElem := rt

	if rt.Kind() != reflect.Ptr {
		return errors.New("argument 2 should be map or ptr")
	}

	rt = rt.(reflect2.PtrType).Elem()
	rtElem = rt

	if rtElem.Kind() != reflect.Struct {
		return errors.New("non-structure type not supported yet")
	}

	s := rtElem.(reflect2.StructType)

	for i := 0; i < s.NumField(); i++ {

		f := s.Field(i)

		min := int64(0)
		max := int64(math.MaxInt64)
		name := f.Tag().Get("name")
		rule := f.Tag().Get("rule")
		msg := f.Tag().Get("msg")
		required := f.Tag().Get("required")
		def := f.Tag().Get("default")
		nums := len(def)

		if len(name) == 0 {
			name = strings.ToLower(f.Name())
		}

		if v, err := strconv.ParseInt(f.Tag().Get("min"), 10, 64); err == nil {
			min = v
		}
		if v, err := strconv.ParseInt(f.Tag().Get("max"), 10, 64); err == nil {
			max = v
		}

		defaultVal := ""
		if string(ctx.Method()) == "GET" {
			defaultVal = strings.TrimSpace(string(ctx.QueryArgs().Peek(name)))
		} else if string(ctx.Method()) == "POST" {
			defaultVal = strings.TrimSpace(string(ctx.PostArgs().Peek(name)))
		}

		check := true //默认需要校验
		if defaultVal == "" {
			if nums > 0 {
				defaultVal = def
			}

			// 是必选参数，且没有默认值
			if required != "0" && defaultVal == "" {
				if rule == "none" {
					check = false
				} else {
					return errors.New(name + " not found")
				}
			} else {
				check = false
			}
		}

		if check {
			switch rule {
			case "digit":
				if !CheckStringDigit(defaultVal) || !CheckIntScope(defaultVal, min, max) {
					return errors.New(msg)
				}
			case "digitString":
				if !CheckStringDigit(defaultVal) || !CheckStringLength(defaultVal, int(min), int(max)) {
					return errors.New(msg)
				}
			case "sDigit":
				if !CheckStringCommaDigit(defaultVal) || !CheckStringLength(defaultVal, int(min), int(max)) {
					return errors.New(msg)
				}
			case "sAlpha":
				if !CheckStringCommaAlpha(defaultVal) || !CheckStringLength(defaultVal, int(min), int(max)) {
					return errors.New(msg)
				}
			case "url":
				if !CheckUrl(defaultVal) {
					return errors.New(msg)
				}
			case "alnum":
				if !CheckStringAlnum(defaultVal) || !CheckStringLength(defaultVal, int(min), int(max)) {
					return errors.New(msg)
				}
			case "priv":
				if !isPriv(defaultVal) {
					return errors.New(msg)
				}
			case "dateTime":
				if !CheckDateTime(defaultVal) {
					return errors.New(msg)
				}
			case "date":
				if !CheckDate(defaultVal) {
					return errors.New(msg)
				}
			case "time":
				if !checkTime(defaultVal) {
					return errors.New(msg)
				}
			case "chn":
				if !CheckStringCHN(defaultVal) {
					return errors.New(msg)
				}
			case "module":
				if !CheckStringModule(defaultVal) || !CheckStringLength(defaultVal, int(min), int(max)) {
					return errors.New(msg)
				}
			case "float":
				if !CheckFloat(defaultVal) {
					return errors.New(msg)
				}
			case "vnphone":
				if !IsVietnamesePhone(defaultVal) {
					return errors.New(msg)
				}
			case "filter":
				if !CheckStringLength(defaultVal, int(min), int(max)) {
					return errors.New(msg)
				}

				defaultVal = FilterInjection(defaultVal)
			case "uname": //会员账号
				if !CheckUName(defaultVal, int(min), int(max)) {
					return errors.New(msg)
				}
			case "upwd": //会员密码
				if !CheckUPassword(defaultVal, int(min), int(max)) {
					return errors.New(msg)
				}
			default:
				break
			}
		}

		switch f.Type().Kind() {
		case reflect.Bool:
			if val, err := strconv.ParseBool(defaultVal); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Int:
			if val, err := strconv.Atoi(defaultVal); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Int8:
			if val, err := strconv.ParseInt(defaultVal, 10, 8); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Int16:
			if val, err := strconv.ParseInt(defaultVal, 10, 16); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Int32:
			if val, err := strconv.ParseInt(defaultVal, 10, 32); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Int64:
			if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Uint:
			if val, err := strconv.ParseUint(defaultVal, 10, 64); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Uint8:
			if val, err := strconv.ParseUint(defaultVal, 10, 8); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Uint16:
			if val, err := strconv.ParseUint(defaultVal, 10, 16); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Uint32:
			if val, err := strconv.ParseUint(defaultVal, 10, 32); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Uint64:
			if val, err := strconv.ParseUint(defaultVal, 10, 64); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Uintptr:
			if val, err := strconv.ParseUint(defaultVal, 10, 64); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Float32:
			if val, err := strconv.ParseFloat(defaultVal, 32); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.Float64:
			if val, err := strconv.ParseFloat(defaultVal, 64); err == nil {
				f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(val))
			}
		case reflect.String:
			f.UnsafeSet(reflect2.PtrOf(objs), reflect2.PtrOf(defaultVal))
		}
		//fmt.Println("MatchType = ", f.MatchType().Kind())
	}

	return nil
}

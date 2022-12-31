package middleware

import (
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"merchant/contrib/session"
	"merchant/model"
	"strconv"

	"github.com/valyala/fastjson"

	"github.com/valyala/fasthttp"
)

var allows = map[string]bool{
	"/merchant/tree":               true,
	"/merchant/captcha":            true,
	"/merchant/admin/login":        true,
	"/merchant/version":            true,
	"/merchant/pprof/":             true,
	"/merchant/pprof/block":        true,
	"/merchant/pprof/allocs":       true,
	"/merchant/pprof/cmdline":      true,
	"/merchant/pprof/goroutine":    true,
	"/merchant/pprof/heap":         true,
	"/merchant/pprof/profile":      true,
	"/merchant/pprof/trace":        true,
	"/merchant/pprof/threadcreate": true,
}

var (
	otpList = map[string]bool{
		"/merchant/group/update": true,
	}
)

func CheckTokenMiddleware(ctx *fasthttp.RequestCtx) error {

	path := string(ctx.Path())
	//fmt.Println("path = ", path)
	if _, ok := allows[path]; ok {
		return nil
	}

	data, err := session.Get(ctx)
	if err != nil {
		return errors.New(`{"status":false,"data":"token"}`)
	}

	ctx.SetUserValue("token", data)

	// 需要otp校验的路由
	//if _, ok := otpList[path]; ok {
	//	if !otp(ctx, data) {
	//		return errors.New(`{"status":false,"data":"otp"}`)
	//	}
	//}

	gid := fastjson.GetString(data, "group_id")
	if path != "/merchant/admin/login" && path != "/merchant/admin/logout" {

		permission := model.PrivCheck(path, gid)
		if permission != nil {
			fmt.Println("path = ", path)
			fmt.Println("gid = ", gid)
			fmt.Println("permission = ", permission)
			return errors.New(`{"status":false,"data":"permission denied"}`)
		}
	}

	if path == "/merchant/admin/logout" {
		ctx.SetUserValue("token", data)
		model.AdminLogout(ctx)
		session.Destroy(ctx)
		return errors.New(`{"status":true,"data":"success"}`)
	}

	return nil
}

func otp(ctx *fasthttp.RequestCtx, data []byte) bool {

	seamo := ""
	if ctx.IsPost() {
		seamo = string(ctx.PostArgs().Peek("code"))
	} else if ctx.IsGet() {
		seamo = string(ctx.QueryArgs().Peek("code"))
	} else {
		return false
	}

	key := fastjson.GetString(data, "seamo")

	fmt.Println("seamo= ", seamo)
	fmt.Println("key= ", key)
	slat := helper.TOTP(key, 60)
	if s, err := strconv.Atoi(seamo); err != nil || s != slat {
		fmt.Println(s, slat)
		return false
	}
	//fmt.Println("seamo = ", key)

	return true
}

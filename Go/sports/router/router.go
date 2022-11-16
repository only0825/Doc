package router

import (
	"fmt"
	"runtime/debug"
	"time"

	"goblog/controller"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var (
	ApiTimeoutMsg = `{"status": "false","data":"服务器响应超时，请稍后重试"}`
	ApiTimeout    = time.Second * 30
	router        *fasthttprouter.Router
	buildInfo     BuildInfo
)

type BuildInfo struct {
	GitReversion   string
	BuildTime      string
	BuildGoVersion string
}

func apiServerPanic(ctx *fasthttp.RequestCtx, rcv interface{}) {

	err := rcv.(error)
	fmt.Println(err)
	debug.PrintStack()

	if r := recover(); r != nil {
		fmt.Println("recovered failed", r)
	}

	ctx.SetStatusCode(500)
	return
}

func Version(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("text/html; charset=utf-8")
	fmt.Fprintf(ctx, "reportApi<br />Git reversion = %s<br />Build Time = %s<br />Go version = %s<br />System Time = %s<br />",
		buildInfo.GitReversion, buildInfo.BuildTime, buildInfo.BuildGoVersion, ctx.Time())

	//ctx.Request.Header.VisitAll(func (key, value []byte) {
	//	fmt.Fprintf(ctx, "%s: %s<br/>", string(key), string(value))
	//})
}

// SetupRouter 设置路由列表
func SetupRouter(b BuildInfo) *fasthttprouter.Router {

	router = fasthttprouter.New()
	router.PanicHandler = apiServerPanic

	buildInfo = b

	// 用户服务
	userCtl := new(controller.UserController)

	get("/log/td/version", Version)

	// 注册用户
	get("/user/reg", userCtl.Reg)
	return router
}

// get is a shortcut for router.GET(path string, handle fasthttp.RequestHandler)
func get(path string, handle fasthttp.RequestHandler) {
	router.GET(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// head is a shortcut for router.HEAD(path string, handle fasthttp.RequestHandler)
func head(path string, handle fasthttp.RequestHandler) {
	router.HEAD(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// options is a shortcut for router.OPTIONS(path string, handle fasthttp.RequestHandler)
func options(path string, handle fasthttp.RequestHandler) {
	router.OPTIONS(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// post is a shortcut for router.POST(path string, handle fasthttp.RequestHandler)
func post(path string, handle fasthttp.RequestHandler) {
	router.POST(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// put is a shortcut for router.PUT(path string, handle fasthttp.RequestHandler)
func put(path string, handle fasthttp.RequestHandler) {
	router.PUT(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// patch is a shortcut for router.PATCH(path string, handle fasthttp.RequestHandler)
func patch(path string, handle fasthttp.RequestHandler) {
	router.PATCH(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

// delete is a shortcut for router.DELETE(path string, handle fasthttp.RequestHandler)
func delete(path string, handle fasthttp.RequestHandler) {
	router.DELETE(path, fasthttp.TimeoutHandler(handle, ApiTimeout, ApiTimeoutMsg))
}

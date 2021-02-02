package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"time"
)

// errorCodeMap 定义了程序需要使用的 HTTP 错误表
var errorCodeMap = map[int]string{
	-1:  "`hex` and `url` query param is not set",
	404: "the resource you specific is not exist.",
	405: "The request method you required is not allowed.",
}

// FailWithMsg 作用于 HTTP 处理方法失败的情况（带有指定消息）
func FailWithMsg(ctx *fasthttp.RequestCtx, code int, msg string) {
	ctx.SetContentType("application/json;charset=utf-8")
	if code > 0 {
		ctx.SetStatusCode(code)
	}
	data, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": msg,
		"data":    []int{},
		"ts":      time.Now().UnixNano() / 1e6,
	}) // 原则上内容是固定的，因此忽略错误
	_, _ = ctx.Write(data) // 原则上内容是固定的，因此显式忽略错误
}

// Fail 作用于 HTTP 处理方法失败的情况（不带有指定消息），消息直接使用错误表
func Fail(ctx *fasthttp.RequestCtx, code int) {
	ctx.SetContentType("application/json;charset=utf-8")
	if code > 0 {
		ctx.SetStatusCode(code)
	}
	msg, ok := errorCodeMap[code]
	if !ok {
		msg = "This error code is not described. Please contact author."
	}
	data, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": msg,
		"data":    []int{},
		"ts":      time.Now().UnixNano() / 1e6,
	}) // 原则上内容是固定的，因此忽略错误
	_, _ = ctx.Write(data) // 原则上内容是固定的，因此显式忽略错误
}

// PixivFetcherHandler 用于响应 Pixiv 资源获取的 HTTP 请求
func PixivFetcherHandler(ctx *fasthttp.RequestCtx) {
	// 检查方法
	if !ctx.IsGet() {
		Fail(ctx, 405)
		return
	}
	// 处理 query 参数
	query := ctx.QueryArgs()
	var url string
	if query.Has("url") {
		url = string(query.Peek("url"))
	} else if query.Has("hex") {
		var err error
		if url, err = base64DecodeString(query.Peek("hex")); err != nil {
			FailWithMsg(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		Fail(ctx, -1)
		return
	}
	// 处理反代
	body, mime, err := fetchPixivResources(url)
	if err != nil {
		FailWithMsg(ctx, 500, err.Error())
		return
	}
	ctx.SetContentType(mime)
	defer body.Close()
	_, err = io.Copy(ctx, body)
	if err != nil {
		log.Error(err)
	}
}

// NotFoundHandler 用于响应 404 情况
func NotFoundHandler(ctx *fasthttp.RequestCtx) {
	Fail(ctx, http.StatusNotFound)
}

// Handler 是 fasthttp 的监听回调处理函数
func Handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		PixivFetcherHandler(ctx)
	default:
		NotFoundHandler(ctx)
	}
}

// RunHTTPServer 用于启动 HTTP 服务
func RunHTTPServer() {
	port := viper.GetInt("server.port")
	dst := fmt.Sprintf(":%v", port)
	log.Info("[proxy] HTTP 服务已启动，您可以通过 " + dst + " 访问")
	if err := fasthttp.ListenAndServe(dst, Handler); err != nil {
		log.Fatal(err)
	}
}

package myGin

import (
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 / 2

//封装Request、ResponseWriter
type Context struct{
	Request *http.Request
	//web框架的最重要的就是输出数据给客户端，
	//这里的输出逻辑我们极有可能需要自己封装一些框架自带的方法。
	//所以我们不妨自定义一个结构R，来实现基本的http.ResponseWriter。并且实现一些具体的其他方法。
	writermen R
	//ResponseWriter 包含了：
	// http.ResponseWriter，http.Hijacker，http.Flusher，http.CloseNotifier和额外方法
	// 暴露给handler，是writermen的复制
	Writer    ResponseWriter

	handlers HandlersChain
	index    int8
}

var _ ResponseWriter = &R{}


var _ ResponseWriter = &R{}
var _ http.ResponseWriter = &R{}
//var _ http.Handler = &Context{}

func (c *Context) reset() {
	c.Writer = &c.writermen
	c.handlers = nil
	c.index = -1
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

//通过 c.Render() 这个渲染的通用方法来适配不同的渲染器
func (c *Context) Render(code int, r Render) {
	//c.Writer.WriteHeader(code)

	if !checkStatus(code) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	}

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

// checkStatus is a copy of http.bodyAllowedForStatus non-exported function.
func checkStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}


// String 将给定字符串写入到响应体中。
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Render(code, String{Format: format, Data: values})
}

// JSON 将给定结构序列化为 JSON 到响应主中。
// 将 Content-Type 设置为 “application/json” 。
func (c *Context) JSON(code int, obj interface{}) {
	c.Render(code, JSON{Data: obj})
}
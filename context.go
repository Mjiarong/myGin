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
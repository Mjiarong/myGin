package myGin

import (
	"bufio"
	"net"
	"net/http"
)

//封装http.ResponseWriter接口，该接口在http中被实现response
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// get the http.Pusher for server push
	Pusher() http.Pusher
}

//构建一个实现ResponseWriter的结构体类型
type R struct {
	http.ResponseWriter
	size   int
	status int
}

//通过http.ResponseWriter分别实现了Hijacker、Flusher、CloseNotifier接口，接管http
func (r *R) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	//将http.ResponseWriter断言为http.Hijacker接口并调用其中的Hijack()方法
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

func (r *R) CloseNotify() <-chan bool {
	return r.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (r *R) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

func (r *R) Pusher() (pusher http.Pusher) {
	if pusher, ok := r.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

func (r *R) reset(writer http.ResponseWriter) {
	r.ResponseWriter = writer
}
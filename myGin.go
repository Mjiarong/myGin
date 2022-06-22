package myGin

import (
	"fmt"
	"net/http"
	"sync"
)

type Engine struct {
	pool  sync.Pool
	mu    sync.RWMutex
	maps     methodMaps
	Router
}

var _ IRouter = &Engine{}

type HandlersChain []handler

type handler func(*Context)

func NewEngine() *Engine {
	engine:= new(Engine)
	engine.Router.engine = engine
	engine.pool.New = func() interface{} {
		return engine.newContext()
	}
	return engine
}

func (engine *Engine) Run(addr string, handler http.Handler)  {
	if err:=http.ListenAndServe(addr,handler);err!=nil{
		fmt.Println("start http server fail:",err)
	}
	return
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	mp:= engine.maps.get(method)
	if mp == nil {
		mp = make(map[string]HandlersChain)
		engine.maps = append(engine.maps, methodMap{method: method, hmap: mp})
	}
	mp.addRoute(path, handlers)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := engine.pool.Get().(*Context)

	c.writermen.reset(w)
	c.Request = r
	c.reset()

	engine.handleHTTPRequest(c)
	engine.pool.Put(c)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	// Find methodMap of the tree for the given HTTP method
	m := engine.maps
	for i, ml := 0, len(m); i < ml; i++ {
		if m[i].method != httpMethod {
			continue
		}
		hmap := m[i].hmap
		handler,ok:= hmap[rPath]
		if ok {
			c.handlers = handler
			c.Next()
			return
		}
	}
}

func (engine *Engine) newContext() *Context {
	return &Context{}
}

func (engine *Engine) Use(middleware ...handler) IRouter {
	engine.Router.Use(middleware...)
	return engine
}

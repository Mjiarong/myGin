package myGin

import "net/http"


type Router struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

var _ IRouter = &Router{}

// IRouter defines all router handle interface.
type IRouter interface {
	Use(...handler) IRouter

	//Handle(string, string, ...handler) IRouter
	//Any(string, ...handler) IRouter
	GET(string, ...handler) IRouter
	POST(string, ...handler) IRouter
	DELETE(string, ...handler) IRouter
	PATCH(string, ...handler) IRouter
	PUT(string, ...handler) IRouter
	OPTIONS(string, ...handler) IRouter
	HEAD(string, ...handler) IRouter

	//StaticFile(string, string) IRoutes
	//Static(string, string) IRoutes
	//StaticFS(string, http.FileSystem) IRoutes
}



func (group *Router) handle(httpMethod, relativePath string, handlers HandlersChain) IRouter {
	fullPath := group.combineAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(httpMethod, fullPath, handlers)
	return group.returnObj()
}

func (group *Router) POST(relativePath string, handlers ...handler) IRouter {
	return group.handle(http.MethodPost, relativePath, handlers)
}

func (group *Router) GET(relativePath string, handlers ...handler) IRouter {
	return group.handle(http.MethodGet, relativePath, handlers)
}

func (group *Router) DELETE(relativePath string, handlers ...handler) IRouter {
	return group.handle(http.MethodDelete, relativePath, handlers)
}

func (group *Router) PATCH(relativePath string, handlers ...handler) IRouter {
	return group.handle(http.MethodPatch, relativePath, handlers)
}

func (group *Router) PUT(relativePath string, handlers ...handler) IRouter {
	return group.handle(http.MethodPut, relativePath, handlers)
}

func (group *Router) OPTIONS(relativePath string, handlers ...handler) IRouter {
	return group.handle(http.MethodOptions, relativePath, handlers)
}

func (group *Router) HEAD(relativePath string, handlers ...handler) IRouter {
	return group.handle(http.MethodHead, relativePath, handlers)
}

func (group *Router) returnObj() IRouter {
	return group.engine
}

func (group *Router) Use(middleware ...handler) IRouter {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

func (group *Router) Group(relativePath string, handlers ...handler) *Router {
	return &Router{
		Handlers: group.combineHandlers(handlers),
		basePath: group.combineAbsolutePath(relativePath),
		engine:   group.engine,
	}
}

//将定义的公用中间件和路由相关的中间件合并
func (group *Router) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

//将路由的相对路径和路由组的路径合并
func (group *Router) combineAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}
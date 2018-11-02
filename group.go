package autodoc

import "github.com/labstack/echo"

// 由于 Group 方法返回的类型和调用者相同，只能将所有方法提出到一个 interface
type GroupInterfaceFull interface {
	Use(m ...echo.MiddlewareFunc)
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Any(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Match(methods []string, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) []*echo.Route
	Group(prefix string, m ...echo.MiddlewareFunc) GroupInterface
	Static(prefix, root string)
	File(path, file string)
	Add(method, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type GroupInterface interface {
	Group(prefix string, middleware ...echo.MiddlewareFunc) GroupInterface
	Use(middleware ...echo.MiddlewareFunc)
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type _group struct {
	group *echo.Group
}

// handler.go 中调用
func NewGroup(e *echo.Echo) GroupInterface {
	return &_group{e.Group("")}
}

func (g *_group) Group(prefix string, middleware ...echo.MiddlewareFunc) GroupInterface {
	return &_group{g.group.Group(prefix, middleware...)}
}

func (g *_group) Use(middleware ...echo.MiddlewareFunc) {
	g.group.Use(middleware...)
}

func (g *_group) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return g.group.GET(path, h, m...)
}

func (g *_group) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return g.group.POST(path, h, m...)
}

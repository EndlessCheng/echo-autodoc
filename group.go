package autodoc

import "github.com/labstack/echo"

type GroupInterface interface {
	Group(prefix string, middleware ...echo.MiddlewareFunc) GroupInterface
	Use(middleware ...echo.MiddlewareFunc)
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type autoDocGroup struct {
	group *echo.Group
}

func NewAutoDocGroup(e *echo.Echo) GroupInterface {
	return &autoDocGroup{e.Group("")}
}

func (mg *autoDocGroup) Group(prefix string, middleware ...echo.MiddlewareFunc) GroupInterface {
	return &autoDocGroup{mg.group.Group(prefix, middleware...)}
}

func (mg *autoDocGroup) Use(middleware ...echo.MiddlewareFunc) {
	mg.group.Use(middleware...)
}

func (mg *autoDocGroup) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return mg.group.GET(path, h, m...)
}

func (mg *autoDocGroup) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return mg.group.POST(path, h, m...)
}

package autodoc

import (
	"strings"

	"github.com/labstack/echo"
)

func getRealHandlerName(name string) string {
	const _suffix = ")-fm"
	if strings.HasSuffix(name, _suffix) {
		splits := strings.Split(name, ".")
		last := splits[len(splits)-1]
		name = last[:len(last)-len(_suffix)]
	}
	return strings.Title(name)
}

type APICollector struct {
	prefix string
	group  *echo.Group
}

func NewAPICollector() *APICollector {
	return &APICollector{"", echo.New().Group("")}
}

func (g *APICollector) Group(prefix string, middleware ...echo.MiddlewareFunc) GroupInterface {
	return &APICollector{g.prefix + prefix, g.group.Group(prefix, middleware...)}
}

func (g *APICollector) Use(middleware ...echo.MiddlewareFunc) {
	g.group.Use(middleware...)
}

// TODO: run middleware
func (g *APICollector) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	r := g.group.GET(path, h, m...)
	docGen.add(getRealHandlerName(r.Name), r.Method, r.Path)
	h(docGen)
	return r
}

// TODO: run middleware
func (g *APICollector) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	r := g.group.POST(path, h, m...)
	docGen.add(getRealHandlerName(r.Name), r.Method, r.Path)
	h(docGen)
	return r
}

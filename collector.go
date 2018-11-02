package autodoc

import (
	"strings"

	"github.com/labstack/echo"
	"io/ioutil"
	"os"
	"runtime"
)

type apiCollector struct {
	prefix string
	group  *echo.Group
}

// handler_test.go 中调用
func NewAPICollector() *apiCollector {
	return &apiCollector{"", echo.New().Group("")}
}

func (g *apiCollector) Group(prefix string, middleware ...echo.MiddlewareFunc) GroupInterface {
	return &apiCollector{g.prefix + prefix, g.group.Group(prefix, middleware...)}
}

func (g *apiCollector) Use(middleware ...echo.MiddlewareFunc) {
	g.group.Use(middleware...)
}

func getRealHandlerName(name string) string {
	const _suffix = ")-fm"
	if strings.HasSuffix(name, _suffix) {
		splits := strings.Split(name, ".")
		last := splits[len(splits)-1]
		name = last[:len(last)-len(_suffix)]
	}
	return name
}

// TODO: run middleware
// 考虑到不同的 GET/POST 可能会调到同一个 handler，在 GET/POST 处提取注释是最准确的
func (g *apiCollector) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	r := g.group.GET(path, h, m...)

	// 提取注释
	_, filePath, lineno, _ := runtime.Caller(1)
	comments := readAboveComments(filePath, lineno-1)

	// 如果开头包含 [skip gen] 则忽略
	if len(comments) > 0 && comments[0] == SkipGen {
		return r
	}

	docGen.add(getRealHandlerName(r.Name), r.Method, r.Path)
	h(docGen)

	if len(comments) > 0 {
		docGen.currentAPI().title = comments[0]
		docGen.currentAPI().description = strings.Join(comments[1:], "\n")
	}

	return r
}

func (g *apiCollector) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	r := g.group.POST(path, h, m...)

	// 提取注释
	_, filePath, lineno, _ := runtime.Caller(1)
	comments := readAboveComments(filePath, lineno-1)

	// 如果开头包含 [skip gen] 则忽略
	if len(comments) > 0 && comments[0] == SkipGen {
		return r
	}

	docGen.add(getRealHandlerName(r.Name), r.Method, r.Path)
	h(docGen)

	if len(comments) > 0 {
		docGen.currentAPI().title = comments[0]
		docGen.currentAPI().description = strings.Join(comments[1:], "\n")
	}

	return r
}

func (g *apiCollector) GenerateMarkdown() string {
	return docGen.generateMarkdown()
}

func (g *apiCollector) GenerateMarkdownToFile(filePath string) error {
	return ioutil.WriteFile(filePath, []byte(g.GenerateMarkdown()), os.ModePerm)
}

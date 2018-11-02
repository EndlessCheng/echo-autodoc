package example

import (
	"testing"
	"github.com/EndlessCheng/echo-autodoc"
)

func Test_setHandlers(t *testing.T) {
	// 初始化默认值
	autodoc.SetQueryParams(autodoc.Param{Type: "string", Name: "isbn", Desc: "ISBN"})

	g := autodoc.NewAPICollector()
	setHandlers(g)

	t.Log(g.GenerateMarkdown())
}

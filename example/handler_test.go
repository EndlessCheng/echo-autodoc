package example

import (
	"testing"
	"github.com/EndlessCheng/echo-autodoc"
)

func Test_setHandlers(t *testing.T) {
	// （可选）设置 URL 默认值
	autodoc.SetQueryParams(autodoc.Param{Type: "string", Name: "isbn", Desc: "ISBN"})

	g := autodoc.NewAPICollector()
	setHandlers(g)

	if err := g.GenerateMarkdownToFile("README.md"); err != nil {
		t.Fatal(err)
	}
}

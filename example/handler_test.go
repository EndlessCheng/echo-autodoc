package example

import (
	"testing"
	"github.com/EndlessCheng/echo-autodoc"
)

func Test_setHandlers(t *testing.T) {
	g := autodoc.NewAPICollector()
	setHTTPHandler(g)
	if err := g.GenerateMarkdownToFile("README.md"); err != nil {
		t.Fatal(err)
	}
}

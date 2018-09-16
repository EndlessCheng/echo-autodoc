package autodoc

import (
	"testing"
	"net/http"
)

type Foo struct {
	Title string `json:"title" desc:"标题"`
}

type Fuu struct {
	Foos []Foo `json:"foos" desc:"Foos"`
}

func TestDocGenerator_Bind(t *testing.T) {
	docGen.add("test", "POST", "/")

	d := struct {
		Name string `json:"name" desc:"姓名"`
		Age  int    `json:"name" desc:"年龄"`
		Male bool   `json:"male" desc:"是否位男性"`
	}{}
	docGen.Bind(&d)

	docGen.JSON(http.StatusOK, &Fuu{[]Foo{{"123"}}})

	t.Log(docGen.generateMarkdown())
}
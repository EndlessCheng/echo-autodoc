package autodoc

import (
	"fmt"
	"strings"
)

type Param struct {
	Type string
	Name string
	Desc string
}

type paramList []Param

func (pl paramList) String() string {
	sb := stringBuilder{}
	sb.WriteString(`|参数|描述|取值|
|-------|--------|--------|
`)
	for _, p := range pl {
		sb.WriteString(fmt.Sprintf("|%s|%s|%s|\n", p.Name, p.Desc, p.Type))
	}
	sb.WriteString("\n")
	return sb.String()
}

type API struct {
	handlerName         string
	method              string
	path                string
	queryParams         paramList
	jsonParams          paramList
	formParams          paramList
	responseParams      paramList
	responseExampleJSON string
}

func (a *API) AddQueryParam(p Param) {
	a.queryParams = append(a.queryParams, p)
}

func (a *API) AddJSONParam(p Param) {
	a.jsonParams = append(a.jsonParams, p)
}

func (a *API) AddFormParam(p Param) {
	a.jsonParams = append(a.jsonParams, p)
}

func (a *API) AddResponseParam(p Param) {
	a.responseParams = append(a.responseParams, p)
}

func (a *API) String() string {
	sb := stringBuilder{}
	// TODO: find DOC
	sb.WriteString(fmt.Sprintf("### %s\n\n`%s %s`\n", strings.Title(a.handlerName), a.method, a.path))

	if len(a.queryParams) > 0 {
		sb.WriteString("\nURL 参数\n\n")
		sb.WriteString(a.queryParams.String())
	}

	if len(a.jsonParams) > 0 {
		sb.WriteString("\nJSON 参数\n\n")
		sb.WriteString(a.jsonParams.String())
	}

	if len(a.formParams) > 0 {
		sb.WriteString("\n表单参数\n\n")
		sb.WriteString(a.formParams.String())
	}

	if len(a.responseParams) > 0 {
		sb.WriteString("\n返回\n\n")
		sb.WriteString(a.responseParams.String())
	}

	if a.responseExampleJSON != "" {
		sb.WriteString(fmt.Sprintf("例如：\n```json\n%s\n```\n", string(a.responseExampleJSON)))
	}

	return sb.String()
}

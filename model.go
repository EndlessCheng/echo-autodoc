package autodoc

import "fmt"

type param struct {
	type_ string
	name  string
	desc  string
}

type paramList []param

func (pl paramList) String() string {
	sb := stringBuilder{}
	sb.WriteString(`|参数|描述|取值|
|-------|--------|--------|
`)
	for _, p := range pl {
		sb.WriteString(fmt.Sprintf("|%s|%s|%s|\n", p.name, p.desc, p.type_))
	}
	sb.WriteString("\n")
	return sb.String()
}

type api struct {
	handlerName       string
	method            string
	path              string
	queryParams       paramList
	jsonParams        paramList
	formParams        paramList
	returnParams      paramList
	returnExampleJSON string
}

func (a *api) addQueryParam(type_ string, name string, desc string) {
	a.queryParams = append(a.queryParams, param{type_, name, desc})
}

func (a *api) addJSONParam(type_ string, name string, desc string) {
	a.jsonParams = append(a.jsonParams, param{type_, name, desc})
}

func (a *api) addFormParam(type_ string, name string, desc string) {
	a.jsonParams = append(a.jsonParams, param{type_, name, desc})
}

func (a *api) addReturnParam(type_ string, name string, desc string) {
	a.returnParams = append(a.returnParams, param{type_, name, desc})
}

func (a *api) String() string {
	sb := stringBuilder{}
	// TODO: find DOC
	sb.WriteString(fmt.Sprintf("### %s\n\n`%s %s`\n", a.handlerName, a.method, a.path))

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

	if len(a.returnParams) > 0 {
		sb.WriteString("\n返回\n\n")
		sb.WriteString(a.returnParams.String())
	}

	if a.returnExampleJSON != "" {
		sb.WriteString(fmt.Sprintf("例如：\n```json\n%s\n```\n", string(a.returnExampleJSON)))
	}

	return sb.String()
}

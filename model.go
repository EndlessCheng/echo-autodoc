package autodoc

import (
	"fmt"
	"strings"
	log "github.com/sirupsen/logrus"
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
	return sb.String()
}

type api struct {
	title       string // 从注释中提取
	description string // 从注释中提取

	handlerName         string
	method              string
	path                string
	queryParams         paramList
	jsonParams          paramList
	formParams          paramList
	responseParams      paramList
	responseExampleJSON string
}

func (a *api) addQueryParam(p ...Param) {
	a.queryParams = append(a.queryParams, p...)
}

func (a *api) addJSONParam(p ...Param) {
	a.jsonParams = append(a.jsonParams, p...)
}

func (a *api) addFormParam(p ...Param) {
	a.formParams = append(a.formParams, p...)
}

func (a *api) addResponseParam(p ...Param) {
	a.responseParams = append(a.responseParams, p...)
}

func (a *api) String() string {
	sb := stringBuilder{}

	var apiTitle string
	if a.title != "" {
		apiTitle = a.title
	} else {
		apiTitle = strings.Title(a.handlerName)
	}
	sb.WriteString(fmt.Sprintf("\n### %s\n", apiTitle))

	sb.WriteString(fmt.Sprintf("\n`%s %s`\n", a.method, a.path))

	if a.description != "" {
		sb.WriteString("\n" + a.description + "\n")
	}

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
		sb.WriteString("\n返回字段\n\n")
		sb.WriteString(a.responseParams.String())
	}

	if a.responseExampleJSON != "" {
		sb.WriteString(fmt.Sprintf("\n返回示例：\n```json\n%s\n```\n", string(a.responseExampleJSON)))
	}

	return sb.String()
}

func (a *api) warnMissingFields() {
	for _, p := range a.queryParams {
		if p.Desc == "" {
			log.Warnf("[%s %s] 缺少 URL 参数描述 - %s", a.method, a.path, p.Name)
		}
	}

	for _, p := range a.jsonParams {
		if p.Desc == "" {
			log.Warnf("[%s %s] 缺少 JSON 参数描述 - %s", a.method, a.path, p.Name)
		}
	}

	for _, p := range a.formParams {
		if p.Desc == "" {
			log.Warnf("[%s %s] 缺少表单参数描述 - %s", a.method, a.path, p.Name)
		}
	}

	for _, p := range a.responseParams {
		if p.Desc == "" {
			log.Warnf("[%s %s] 缺少返回字段描述 - %s", a.method, a.path, p.Name)
		}
	}
}

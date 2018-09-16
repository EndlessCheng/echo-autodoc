package autodoc

var DefaultMarkdownHeader = `# 接口文档

## HTTP 接口
`

var DefaultGetReturn interface{} = 1
var DefaultQueryParamReturn = "1"

var customURLParams []Param
var customJSONParams []Param
var customResponseParams []Param

func SetURLParams(params ...Param) {
	customURLParams = params
}

func SetJSONParams(params ...Param) {
	customJSONParams = params
}

func SetResponseParams(params ...Param) {
	customResponseParams = params
}

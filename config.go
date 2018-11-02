package autodoc

var (
	DefaultMarkdownHeader = `# 接口文档

## HTTP 接口
`
	DefaultGetReturn        interface{} = 1
	DefaultQueryParamReturn             = "1"
)

var (
	customQueryParams    = map[string]Param{}
	customPostJSONParams = map[string]Param{}
	customResponseParams = map[string]Param{}
)

func SetQueryParams(params ...Param) {
	for _, p := range params {
		customQueryParams[p.Name] = p
	}
}

func SetPostJSONParams(params ...Param) {
	for _, p := range params {
		customPostJSONParams[p.Name] = p
	}
}

func SetResponseParams(params ...Param) {
	for _, p := range params {
		customResponseParams[p.Name] = p
	}
}

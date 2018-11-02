package autodoc

var (
	SkipGen = "[skip gen]"

	DefaultMarkdownHeader = `# 接口文档

## HTTP 接口`
	DefaultGetReturn        interface{} = 1
	DefaultQueryParamReturn             = "1"
	DefaultFormValueReturn              = ""
	DefaultFormFileDesc                 = "上传的文件"
)

var (
	customContextGetParams   = map[string]Param{}
	customQueryParams        = map[string]Param{}
	customFormParams         = map[string]Param{}
	customPostJSONParams     = map[string]Param{}
	customResponseJSONParams = map[string]Param{}
)

func SetContextGetParams(params ...Param) {
	for _, p := range params {
		customContextGetParams[p.Name] = p
	}
}

func SetQueryParams(params ...Param) {
	for _, p := range params {
		customQueryParams[p.Name] = p
	}
}

func SetFormParams(params ...Param) {
	for _, p := range params {
		customFormParams[p.Name] = p
	}
}

func SetPostJSONParams(params ...Param) {
	for _, p := range params {
		customPostJSONParams[p.Name] = p
	}
}

func SetResponseJSONParams(params ...Param) {
	for _, p := range params {
		customResponseJSONParams[p.Name] = p
	}
}

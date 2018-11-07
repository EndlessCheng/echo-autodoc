package autodoc

var (
	SkipGen = "[skip gen]"

	DefaultMarkdownHeader = `# 接口文档

## HTTP 接口`
	DefaultMarkdownFooter               = ""
	DefaultGetReturn        interface{} = 1
	DefaultQueryParamReturn             = "1"
	DefaultFormValueReturn              = ""
	DefaultFormFileDesc                 = "上传的文件"
)

var (
	customContextGetParams    = map[string]Param{}
	customQueryParams         = map[string]Param{}
	customFormParams          = map[string]Param{}
	customPostJSONParams      = map[string]Param{}
	customResponseJSONParams  = map[string]Param{}
	globalResponseJSONParams  []Param
	ignoredResponseJSONParams = map[string]Param{}
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

// 同时会设置 ignoredResponseJSONParams
func SetGloablResponseJSONParams(params ...Param) {
	globalResponseJSONParams = append(globalResponseJSONParams, params...)
	for _, p := range params {
		ignoredResponseJSONParams[p.Name] = p
	}
}

func SetIgnoredResponseJSONParams(params ...Param) {
	for _, p := range params {
		ignoredResponseJSONParams[p.Name] = p
	}
}

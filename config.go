package autodoc

import (
	"mime/multipart"
	"reflect"
	"unsafe"
	"net/http"
)

var (
	// c.JSON 时，若传入的数据有零值，可以将其递归地填充为非零值（数值填充成 1，字符串填充成 "1"，布尔值依然为 false）
	FillZeroValue = false

	SkipGen = "[skip gen]"

	DefaultMarkdownHeader = `# 接口文档

## HTTP 接口`
	DefaultMarkdownFooter = ""

	DefaultCookie = http.Cookie{Value: "1"}

	DefaultGetReturn        interface{} = 1
	DefaultQueryParamReturn             = "1"
	DefaultFormValueReturn              = "1"
	DefaultFormFileDesc                 = "上传的文件"

	DefaultMultipartFileHeader        = multipart.FileHeader{Filename: "example.txt", Size: 1}
	DefaultMultipartFileHeaderContent = []byte("A")
)

func init() {
	ModifyFileHeaderContent(&DefaultMultipartFileHeader, DefaultMultipartFileHeaderContent)
}

func ModifyFileHeaderContent(fh *multipart.FileHeader, content []byte) {
	fhVal := reflect.Indirect(reflect.ValueOf(fh))
	ptrToContent := (*[]byte)(unsafe.Pointer(fhVal.FieldByName("content").UnsafeAddr()))
	*ptrToContent = content
}

//

var (
	customGetReturnMap        = map[string]interface{}{}
	customQueryParamReturnMap = map[string]string{}
	customFormValueReturnMap  = map[string]string{}
)

func AddCustomGetReturn(key string, ret interface{}) {
	customGetReturnMap[key] = ret
}

func AddCustomQueryParamReturn(name string, ret string) {
	customQueryParamReturnMap[name] = ret
}

func AddCustomFormValueReturn(name string, ret string) {
	customFormValueReturnMap[name] = ret
}

//

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

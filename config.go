package autodoc

import (
	"mime/multipart"
	"reflect"
	"unsafe"
	"net/http"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strings"
)

// 杂项配置
var (
	// c.JSON 时，若传入的数据有零值，可以将其递归地填充为非零值（数值填充成 1，字符串填充成 "1"，布尔值依然为 false）
	FillZeroValue = false

	SkipGen = "[skip gen]"

	// 当返回值不为 200 时，打印 warning 信息
	WarningWhenNotStatusOK = true

	// 字段缺少描述时打印 warning 信息
	WarningMissingFields = true
)

// Markdown 相关配置
var (
	DefaultMarkdownHeader = `# 接口文档

## HTTP 接口`
	DefaultMarkdownFooter = ""

	DefaultFormFileDesc = "上传的文件"
)

// 默认 echo 操作返回值
// 若要精确设置，见下方的 AddCustomGetReturn 等
var (
	DefaultCookie = http.Cookie{Value: "1"}

	DefaultGetReturn        interface{} = 1
	DefaultQueryParamReturn             = "1"
	DefaultFormValueReturn              = "1"
	// TODO: DefaultPostJSONFieldValue   interface{} = "1"

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

// 自定义类型识别
var (
	customGoTypeToJSONType = map[string]string{
		"time.Time": "string",
		"HashID":    "string",
		"JSONTime":  "string",
	}
)

func AddCustomGoTypeToJSONType(goType string, jsonType string) {
	customGoTypeToJSONType[goType] = jsonType
}

// 精确设置 echo 操作返回值
var (
	customContextGetReturnMap = map[string]interface{}{}
	customQueryParamReturnMap = map[string]string{}
	customFormValueReturnMap  = map[string]string{}
	// TODO: customPostJSONFieldValueMap = map[string]interface{}{}
)

func AddCustomGetReturn(key string, ret interface{}) {
	customContextGetReturnMap[key] = ret
}
func AddCustomContextGetReturn(key string, ret interface{}) {
	customContextGetReturnMap[key] = ret
}
func AddCustomQueryParamReturn(name string, str string) {
	customQueryParamReturnMap[name] = str
}
func AddCustomFormValueReturn(name string, str string) {
	customFormValueReturnMap[name] = str
}

// 解析参数时，若没有相关注释，使用下面的值（一般用于常见的入参出参）
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

// globalResponseJSONParams 为所有返回的 JSON 都会有的字段
// 同时会设置 ignoredResponseJSONParams
func SetGloablResponseJSONParams(params ...Param) {
	globalResponseJSONParams = append(globalResponseJSONParams, params...)
	for _, p := range params {
		ignoredResponseJSONParams[p.Name] = p
	}
}

// 解析返回的 JSON 时，忽略这些字段
func SetIgnoredResponseJSONParams(params ...Param) {
	for _, p := range params {
		ignoredResponseJSONParams[p.Name] = p
	}
}

//

type ContextJSON interface {
	BeforeJSON(code int, i interface{})
	AfterJSON(code int, i interface{})
}

var (
	ContextJSONer ContextJSON = ErrorCodeContextJSONer
)

var (
	EmptyContextJSONer     = &emptyContextJSON{}
	ErrorCodeContextJSONer = &errorCodeContextJSON{1000}
)

type emptyContextJSON struct{}

func (*emptyContextJSON) BeforeJSON(code int, i interface{}) {}
func (*emptyContextJSON) AfterJSON(code int, i interface{})  {}

type errorCodeContextJSON struct {
	errorCodeOK int
}

func (cj *errorCodeContextJSON) BeforeJSON(code int, i interface{}) {
	if code == http.StatusOK {
		data, err := json.Marshal(i)
		if err != nil {
			return
		}

		if !strings.Contains(string(data), `"errcode"`) {
			// log.Warn()
			return
		}

		d := struct {
			ErrCode int `json:"errcode"`
		}{}
		if err := json.Unmarshal(data, &d); err != nil {
			log.WithError(err).Errorln("[c.JSON.BeforeJSON.json.Unmarshal]")
			return
		}

		if d.ErrCode != cj.errorCodeOK {
			log.WithError(err).Errorf("[c.JSON.BeforeJSON] d.ErrCode != cj.errorCodeOK (%d != %d)", d.ErrCode, cj.errorCodeOK)
		}
	}
}
func (*errorCodeContextJSON) AfterJSON(code int, i interface{}) {}

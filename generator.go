package autodoc

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"github.com/labstack/echo"
	"strings"
	"time"
	"runtime"
	log "github.com/sirupsen/logrus"
)

var timeType = reflect.TypeOf(time.Time{})

// docGenerator 实现 echo.Context 接口，从而 hack 进业务逻辑代码（handler.go）中，详细见 collector.go
type docGenerator struct {
	apiList []*api
}

var docGen = &docGenerator{}

func (dg *docGenerator) currentAPI() *api {
	return dg.apiList[len(dg.apiList)-1]
}

func (dg *docGenerator) add(handlerName string, method string, path string) {
	dg.apiList = append(dg.apiList, &api{handlerName: handlerName, method: method, path: path})
}

func (dg *docGenerator) generateMarkdown() string {
	sb := stringBuilder{}
	sb.WriteString(DefaultMarkdownHeader)
	for _, api := range dg.apiList {
		sb.WriteString(fmt.Sprintf("\n%s", api.String()))
	}
	sb.WriteString(DefaultMarkdownFooter)
	return sb.String()
}

func (dg *docGenerator) Request() *http.Request                                  { return nil }
func (dg *docGenerator) SetRequest(r *http.Request)                              {}
func (dg *docGenerator) Response() *echo.Response                                { return nil }
func (dg *docGenerator) IsTLS() bool                                             { return false }
func (dg *docGenerator) IsWebSocket() bool                                       { return false }
func (dg *docGenerator) Scheme() string                                          { return "" }
func (dg *docGenerator) RealIP() string                                          { return "" }
func (dg *docGenerator) Path() string                                            { return "" }
func (dg *docGenerator) SetPath(p string)                                        {}
func (dg *docGenerator) Param(name string) string                                { return "" }
func (dg *docGenerator) ParamNames() []string                                    { return nil }
func (dg *docGenerator) SetParamNames(names ...string)                           {}
func (dg *docGenerator) ParamValues() []string                                   { return nil }
func (dg *docGenerator) SetParamValues(values ...string)                         {}
func (dg *docGenerator) QueryParams() url.Values                                 { return nil }
func (dg *docGenerator) QueryString() string                                     { return "" }
func (dg *docGenerator) FormParams() (url.Values, error)                         { return nil, nil }
func (dg *docGenerator) MultipartForm() (*multipart.Form, error)                 { return nil, nil }
func (dg *docGenerator) SetCookie(cookie *http.Cookie)                           {}
func (dg *docGenerator) Cookies() []*http.Cookie                                 { return nil }
func (dg *docGenerator) Set(key string, val interface{})                         {}
func (dg *docGenerator) Validate(i interface{}) error                            { return nil }
func (dg *docGenerator) Render(code int, name string, data interface{}) error    { return nil }
func (dg *docGenerator) HTML(code int, html string) error                        { return nil }
func (dg *docGenerator) HTMLBlob(code int, b []byte) error                       { return nil }
func (dg *docGenerator) JSONPretty(code int, i interface{}, indent string) error { return nil }
func (dg *docGenerator) JSONBlob(code int, b []byte) error                       { return nil }
func (dg *docGenerator) JSONP(code int, callback string, i interface{}) error    { return nil }
func (dg *docGenerator) JSONPBlob(code int, callback string, b []byte) error     { return nil }
func (dg *docGenerator) XML(code int, i interface{}) error                       { return nil }
func (dg *docGenerator) XMLPretty(code int, i interface{}, indent string) error  { return nil }
func (dg *docGenerator) XMLBlob(code int, b []byte) error                        { return nil }
func (dg *docGenerator) Blob(code int, contentType string, b []byte) error       { return nil }
func (dg *docGenerator) Stream(code int, contentType string, r io.Reader) error  { return nil }
func (dg *docGenerator) File(file string) error                                  { return nil }
func (dg *docGenerator) Attachment(file string, name string) error               { return nil }
func (dg *docGenerator) Inline(file string, name string) error                   { return nil }
func (dg *docGenerator) Redirect(code int, url string) error                     { return nil }
func (dg *docGenerator) Error(err error)                                         {}
func (dg *docGenerator) Handler() echo.HandlerFunc                               { return nil }
func (dg *docGenerator) SetHandler(h echo.HandlerFunc)                           {}
func (dg *docGenerator) Logger() echo.Logger                                     { return nil }
func (dg *docGenerator) Echo() *echo.Echo                                        { return nil }
func (dg *docGenerator) Reset(r *http.Request, w http.ResponseWriter)            {}

func (dg *docGenerator) NoContent(code int) error {
	if WarningWhenNotStatusOK && code != http.StatusOK {
		log.Warnf("[c.NoContent] code is %d", code)
	}
	return nil
}

func (dg *docGenerator) String(code int, s string) error {
	if WarningWhenNotStatusOK && code != http.StatusOK {
		log.Warnf("[c.String] code is %d", code)
	}
	return nil
}

func (dg *docGenerator) Cookie(name string) (*http.Cookie, error) {
	return &DefaultCookie, nil
}

// 显示在 README 中的类型
var validTypes = [...]string{"int", "float", "string", "bool"}

func (dg *docGenerator) isValidType(type_ string) bool {
	for _, t := range validTypes {
		if t == type_ {
			return true
		}
	}
	return false
}

// 类型, 描述
func (dg *docGenerator) parseTailComment(comment string, defaultType string) (type_ string, desc string) {
	if comment == "" {
		return "", ""
	}

	splits := strings.Split(comment, ",")
	if len(splits) == 1 {
		return defaultType, comment
	}

	type_ = strings.TrimSpace(splits[0])
	if !dg.isValidType(type_) {
		return defaultType, comment
	}

	return type_, strings.TrimSpace(strings.Join(splits[1:], ","))
}

// 暂时不加 parseTailComment
func (dg *docGenerator) Get(key string) interface{} {
	if param, ok := customContextGetParams[key]; ok {
		dg.currentAPI().addQueryParam(param)
	}

	if ret, ok := customGetReturnMap[key]; ok {
		return ret
	}
	return DefaultGetReturn
}

func (dg *docGenerator) QueryParam(name string) string {
	_, filePath, lineno, _ := runtime.Caller(1) // skip 的值取决于这行代码离要提取的注释相隔几层调用
	comment := readTailComment(filePath, lineno)

	if type_, desc := dg.parseTailComment(comment, "string"); desc != "" {
		// 特殊大于全局
		dg.currentAPI().addQueryParam(Param{type_, name, desc})
	} else {
		param, ok := customQueryParams[name]
		if !ok {
			param = Param{"string", name, ""}
		}
		dg.currentAPI().addQueryParam(param)
	}

	if ret, ok := customQueryParamReturnMap[name]; ok {
		return ret
	}
	return DefaultQueryParamReturn
}

func (dg *docGenerator) FormFile(name string) (*multipart.FileHeader, error) {
	_, filePath, lineno, _ := runtime.Caller(1) // skip 的值取决于这行代码离要提取的注释相隔几层调用
	comment := readTailComment(filePath, lineno)

	var param Param
	if type_, desc := dg.parseTailComment(comment, "file"); desc != "" {
		param = Param{type_, name, desc}
	} else {
		param = Param{"file", name, DefaultFormFileDesc}
	}
	dg.currentAPI().addFormParam(param)

	return &DefaultMultipartFileHeader, nil
}

func (dg *docGenerator) FormValue(name string) string {
	_, filePath, lineno, _ := runtime.Caller(1) // skip 的值取决于这行代码离要提取的注释相隔几层调用
	comment := readTailComment(filePath, lineno)

	if type_, desc := dg.parseTailComment(comment, "string"); desc != "" {
		// 特殊大于全局
		dg.currentAPI().addFormParam(Param{type_, name, desc})
	} else {
		param, ok := customFormParams[name]
		if !ok {
			param = Param{"string", name, ""}
		}
		dg.currentAPI().addFormParam(param)
	}

	if ret, ok := customFormValueReturnMap[name]; ok {
		return ret
	}
	return DefaultFormValueReturn
}

func dereference(v reflect.Type) reflect.Type {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func typeToString(val reflect.Type) string {
	v := dereference(val)
	switch v.Kind() {
	case reflect.Invalid:
		panic("[typeToString] 暂不支持 " + val.Name())
	case reflect.Bool:
		return "bool"
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return "int"
	case reflect.Float32,
		reflect.Float64:
		return "float"
	case reflect.String:
		return "string"
	case reflect.Array,
		reflect.Slice:
		return typeToString(v.Elem()) + " array"
	default:
		return "object"
	}
}

func getType(i interface{}) string {
	return typeToString(reflect.TypeOf(i))
}

// struct / struct 指针 / struct slice
func parseStructWithPrefix(prefix string, structType reflect.Type) (params []Param) {
	if structType.Kind() == reflect.Ptr || structType.Kind() == reflect.Slice {
		// array 信息已经在上层提取出来了，这里我们只需要内部的信息
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		return params
	}

	// 提取 struct 内部信息
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		var name string
		if jsonTag == "" {
			// 不写 json tag 的话，就用属性类型代替
			name = field.Name
		} else {
			name = strings.Split(jsonTag, ",")[0]
		}
		name = prefix + name

		fieldType := dereference(field.Type)
		isTimeType := fieldType.ConvertibleTo(timeType)

		type_ := typeToString(fieldType)
		// time.Time 固定成 string
		if isTimeType {
			type_ = "string"
		}

		desc := field.Tag.Get("desc")
		params = append(params, Param{type_, name, desc})

		// skip time.Time
		if !isTimeType {
			switch fieldType.Kind() {
			case reflect.Slice,
				reflect.Ptr,
				reflect.Struct:
				params = append(params, parseStructWithPrefix(name+".", field.Type)...)
			}
		}
	}

	return params
}

// struct 或 struct 指针
func parseStruct(structType reflect.Type) []Param {
	return parseStructWithPrefix("", structType)
}

func parseMap(mp map[string]interface{}) (params []Param) {
	for name, val := range mp {
		valType := reflect.TypeOf(val)
		valType = dereference(valType)

		type_ := typeToString(valType)
		// time.Time 固定成 string
		if valType.Name() == "Time" {
			type_ = "string"
		}

		params = append(params, Param{type_, name, ""})

		// skip time.Time
		if valType.Name() != "Time" {
			switch valType.Kind() {
			case reflect.Slice,
				reflect.Ptr,
				reflect.Struct:
				params = append(params, parseStructWithPrefix(name+".", valType)...)
			}
		}
	}

	return params
}

func (dg *docGenerator) Bind(i interface{}) error {
	// 传进来的一定是个 struct 指针
	params := parseStruct(reflect.TypeOf(i).Elem())
	for _, p := range params {
		if customParam, ok := customPostJSONParams[p.Name]; ok {
			p = customParam
		}
		dg.currentAPI().addJSONParam(p)
	}

	return nil
}

func (dg *docGenerator) JSON(code int, i interface{}) error {
	defer ContextJSONer.AfterJSON(code, i)
	ContextJSONer.BeforeJSON(code, i)

	if WarningWhenNotStatusOK && code != http.StatusOK {
		log.Warnf("[c.JSON] code is %d", code)
	}

	if FillZeroValue {
		// TODO
		//switch val := i.(type) {
		//case map[string]interface{}:
		//	for _, _val := range val {
		//		valType := reflect.TypeOf(_val)
		//		// TODO: map中的，不是指针能修改吗？
		//		if valType.Kind() == reflect.Struct {
		//			FillStruct(&_val)
		//		} else if valType.Kind() == reflect.Ptr && valType.Elem().Kind() == reflect.Struct {
		//			FillStruct(_val)
		//		}
		//	}
		//default:
		//	// 否则是个 struct 或 struct 指针
		//	valType := reflect.TypeOf(i)
		//	if valType.Kind() == reflect.Struct {
		//		FillStruct(&i)
		//	} else if valType.Kind() == reflect.Ptr && valType.Elem().Kind() == reflect.Struct {
		//		FillStruct(i)
		//	}
		//}
	}
	data, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		panic(err)
	}
	dg.currentAPI().responseExampleJSON = string(data)

	// TODO: 合并
	var params []Param
	switch val := i.(type) {
	case map[string]interface{}:
		params = parseMap(val)
	default:
		// 否则是个 struct 或 struct 指针
		params = parseStruct(reflect.TypeOf(val))
	}

	for _, p := range params {
		if _, ok := ignoredResponseJSONParams[p.Name]; ok {
			continue
		}
		if customParam, ok := customResponseJSONParams[p.Name]; ok {
			p = customParam
		}
		dg.currentAPI().addResponseParam(p)
	}

	return nil
}

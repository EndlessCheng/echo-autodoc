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
)

// DocGenerator 实现 echo.Context 接口，从而 hack 进业务逻辑代码（handler.go）中，详细见 collector.go
type DocGenerator struct {
	apiList []*API
}

var docGen = &DocGenerator{}

func (dg *DocGenerator) currentAPI() *API {
	return dg.apiList[len(dg.apiList)-1]
}

func (dg *DocGenerator) add(handlerName string, method string, path string) {
	dg.apiList = append(dg.apiList, &API{handlerName: handlerName, method: method, path: path})
}

func (dg *DocGenerator) generateMarkdown() string {
	sb := stringBuilder{}
	sb.WriteString(DefaultMarkdownHeader)
	for _, api := range dg.apiList {
		sb.WriteString(fmt.Sprintf("\n%s", api.String()))
	}
	return sb.String()
}

func (dg *DocGenerator) Request() *http.Request                                  { return nil }
func (dg *DocGenerator) SetRequest(r *http.Request)                              {}
func (dg *DocGenerator) Response() *echo.Response                                { return nil }
func (dg *DocGenerator) IsTLS() bool                                             { return false }
func (dg *DocGenerator) IsWebSocket() bool                                       { return false }
func (dg *DocGenerator) Scheme() string                                          { return "" }
func (dg *DocGenerator) RealIP() string                                          { return "" }
func (dg *DocGenerator) Path() string                                            { return "" }
func (dg *DocGenerator) SetPath(p string)                                        {}
func (dg *DocGenerator) Param(name string) string                                { return "" }
func (dg *DocGenerator) ParamNames() []string                                    { return nil }
func (dg *DocGenerator) SetParamNames(names ...string)                           {}
func (dg *DocGenerator) ParamValues() []string                                   { return nil }
func (dg *DocGenerator) SetParamValues(values ...string)                         {}
func (dg *DocGenerator) QueryParams() url.Values                                 { return nil }
func (dg *DocGenerator) QueryString() string                                     { return "" }
func (dg *DocGenerator) FormParams() (url.Values, error)                         { return nil, nil }
func (dg *DocGenerator) MultipartForm() (*multipart.Form, error)                 { return nil, nil }
func (dg *DocGenerator) Cookie(name string) (*http.Cookie, error)                { return nil, nil }
func (dg *DocGenerator) SetCookie(cookie *http.Cookie)                           {}
func (dg *DocGenerator) Cookies() []*http.Cookie                                 { return nil }
func (dg *DocGenerator) Set(key string, val interface{})                         {}
func (dg *DocGenerator) Validate(i interface{}) error                            { return nil }
func (dg *DocGenerator) Render(code int, name string, data interface{}) error    { return nil }
func (dg *DocGenerator) HTML(code int, html string) error                        { return nil }
func (dg *DocGenerator) HTMLBlob(code int, b []byte) error                       { return nil }
func (dg *DocGenerator) String(code int, s string) error                         { return nil }
func (dg *DocGenerator) JSONPretty(code int, i interface{}, indent string) error { return nil }
func (dg *DocGenerator) JSONBlob(code int, b []byte) error                       { return nil }
func (dg *DocGenerator) JSONP(code int, callback string, i interface{}) error    { return nil }
func (dg *DocGenerator) JSONPBlob(code int, callback string, b []byte) error     { return nil }
func (dg *DocGenerator) XML(code int, i interface{}) error                       { return nil }
func (dg *DocGenerator) XMLPretty(code int, i interface{}, indent string) error  { return nil }
func (dg *DocGenerator) XMLBlob(code int, b []byte) error                        { return nil }
func (dg *DocGenerator) Blob(code int, contentType string, b []byte) error       { return nil }
func (dg *DocGenerator) Stream(code int, contentType string, r io.Reader) error  { return nil }
func (dg *DocGenerator) File(file string) error                                  { return nil }
func (dg *DocGenerator) Attachment(file string, name string) error               { return nil }
func (dg *DocGenerator) Inline(file string, name string) error                   { return nil }
func (dg *DocGenerator) NoContent(code int) error                                { return nil }
func (dg *DocGenerator) Redirect(code int, url string) error                     { return nil }
func (dg *DocGenerator) Error(err error)                                         {}
func (dg *DocGenerator) Handler() echo.HandlerFunc                               { return nil }
func (dg *DocGenerator) SetHandler(h echo.HandlerFunc)                           {}
func (dg *DocGenerator) Logger() echo.Logger                                     { return nil }
func (dg *DocGenerator) Echo() *echo.Echo                                        { return nil }
func (dg *DocGenerator) Reset(r *http.Request, w http.ResponseWriter)            {}

func (dg *DocGenerator) FormValue(name string) string {
	dg.currentAPI().AddFormParam(Param{"string", name, ""})
	return ""
}

func (dg *DocGenerator) FormFile(name string) (*multipart.FileHeader, error) {
	dg.currentAPI().AddFormParam(Param{"file", name, "上传的文件"})
	return nil, nil
}

func (dg *DocGenerator) Get(key string) interface{} {
	return DefaultGetReturn
}

func (dg *DocGenerator) QueryParam(name string) string {
	param, ok := customQueryParams[name]
	if !ok {
		param = Param{"string", name, ""}
	}
	dg.currentAPI().AddQueryParam(param)

	return DefaultQueryParamReturn
}

func dereference(v reflect.Type) reflect.Type {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func typeToString(v reflect.Type) string {
	v = dereference(v)
	switch v.Kind() {
	case reflect.Invalid:
		panic("[typeToString] 代码有误！")
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

func _parseStruct(prefix string, structType reflect.Type) []Param {
	var params []Param

	if structType.Kind() == reflect.Ptr || structType.Kind() == reflect.Slice {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		return params
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldType := dereference(field.Type)

		type_ := typeToString(fieldType)
		// time.Time 固定成 string
		if fieldType.Name() == "Time" {
			type_ = "string"
		}

		// 不写 json tag 的话，就用属性类型代替
		name := field.Tag.Get("json")
		if name == "" {
			name = field.Name
		}
		name = prefix + name

		desc := field.Tag.Get("desc")
		params = append(params, Param{type_, name, desc})

		// time.Time 固定成 string
		if fieldType.Name() != "Time" {
			switch field.Type.Kind() {
			case reflect.Slice,
				reflect.Ptr,
				reflect.Struct:
				params = append(params, _parseStruct(name+".", field.Type)...)
			}
		}
	}

	return params
}

func parseStruct(structType reflect.Type) []Param {
	return _parseStruct("", structType)
}

func (dg *DocGenerator) Bind(i interface{}) error {
	// 传进来的一定是个 struct 指针
	params := parseStruct(reflect.TypeOf(i).Elem())
	for _, p := range params {
		if customParam, ok := customPostJSONParams[p.Name]; ok {
			p = customParam
		}
		dg.currentAPI().AddJSONParam(p)
	}

	return nil
}

func (dg *DocGenerator) JSON(code int, i interface{}) error {
	if code != http.StatusOK {
		// TODO: ignore?
	}

	data, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		panic(err)
	}
	dg.currentAPI().responseExampleJSON = string(data)

	switch val := i.(type) {
	case map[string]interface{}:
		if len(val) > 0 {
			for key, value := range val {
				switch val := value.(type) {
				case map[string]interface{}:
					for name, v := range val {
						dg.currentAPI().AddResponseParam(Param{getType(v), name, ""})
					}
				default:
					dg.currentAPI().AddResponseParam(Param{getType(val), key, ""})
				}
			}
		}
	default:
		// 否则是个 struct 或 struct 指针
		params := parseStruct(reflect.TypeOf(val))
		for _, p := range params {
			if customParam, ok := customResponseParams[p.Name]; ok {
				p = customParam
			}
			dg.currentAPI().AddResponseParam(p)
		}
	}

	return nil
}

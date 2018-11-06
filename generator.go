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
func (dg *docGenerator) Cookie(name string) (*http.Cookie, error)                { return nil, nil }
func (dg *docGenerator) SetCookie(cookie *http.Cookie)                           {}
func (dg *docGenerator) Cookies() []*http.Cookie                                 { return nil }
func (dg *docGenerator) Set(key string, val interface{})                         {}
func (dg *docGenerator) Validate(i interface{}) error                            { return nil }
func (dg *docGenerator) Render(code int, name string, data interface{}) error    { return nil }
func (dg *docGenerator) HTML(code int, html string) error                        { return nil }
func (dg *docGenerator) HTMLBlob(code int, b []byte) error                       { return nil }
func (dg *docGenerator) String(code int, s string) error                         { return nil }
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
func (dg *docGenerator) NoContent(code int) error                                { return nil }
func (dg *docGenerator) Redirect(code int, url string) error                     { return nil }
func (dg *docGenerator) Error(err error)                                         {}
func (dg *docGenerator) Handler() echo.HandlerFunc                               { return nil }
func (dg *docGenerator) SetHandler(h echo.HandlerFunc)                           {}
func (dg *docGenerator) Logger() echo.Logger                                     { return nil }
func (dg *docGenerator) Echo() *echo.Echo                                        { return nil }
func (dg *docGenerator) Reset(r *http.Request, w http.ResponseWriter)            {}

func (dg *docGenerator) Get(key string) interface{} {
	if param, ok := customContextGetParams[key]; ok {
		dg.currentAPI().addQueryParam(param)
	}

	return DefaultGetReturn
}

func (dg *docGenerator) QueryParam(name string) string {
	param, ok := customQueryParams[name]
	if !ok {
		param = Param{"string", name, ""}
	}
	dg.currentAPI().addQueryParam(param)

	return DefaultQueryParamReturn
}

func (dg *docGenerator) FormFile(name string) (*multipart.FileHeader, error) {
	dg.currentAPI().addFormParam(Param{"file", name, DefaultFormFileDesc})
	return nil, nil
}

func (dg *docGenerator) FormValue(name string) string {
	param, ok := customFormParams[name]
	if !ok {
		param = Param{"string", name, ""}
	}
	dg.currentAPI().addFormParam(param)

	return DefaultFormValueReturn
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

		// skip time.Time
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
		for name, v := range val {
			switch _val := v.(type) {
			case map[string]interface{}:
				for _name, _v := range _val {
					p := Param{getType(_v), _name, ""}
					if customParam, ok := customResponseJSONParams[_name]; ok {
						p = customParam
					}
					dg.currentAPI().addResponseParam(p)
				}
			default:
				p := Param{getType(v), name, ""}
				if customParam, ok := customResponseJSONParams[name]; ok {
					p = customParam
				}
				dg.currentAPI().addResponseParam(p)
			}
		}
	default:
		// 否则是个 struct 或 struct 指针
		params := parseStruct(reflect.TypeOf(val))
		for _, p := range params {
			if customParam, ok := customResponseJSONParams[p.Name]; ok {
				p = customParam
			}
			dg.currentAPI().addResponseParam(p)
		}
	}

	return nil
}

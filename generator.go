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
	"os"
)

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
	sb.WriteString(`# 接口文档

## HTTP 接口
`)
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
func (dg *docGenerator) FormValue(name string) string {
	docGen.currentAPI().addFormParam("string", name, "")
	return ""
}
func (dg *docGenerator) FormFile(name string) (*multipart.FileHeader, error) {
	docGen.currentAPI().addFormParam("file", name, "上传的文件")
	return nil, nil
}

func (dg *docGenerator) Get(key string) interface{} {
	return DefaultGetReturn
}

func (dg *docGenerator) QueryParam(name string) string {
	return DefaultQueryParamReturn
}

func getTypeType(v reflect.Type) string {
	switch v.Kind() {
	case reflect.Invalid:
		panic("[_getValueType] 代码有误！")
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.String:
		return "string"
	case reflect.Slice:
		return getTypeType(v.Elem()) + " array"
	default:
		return "object" // TODO: more detail?
	}
}

func getType(i interface{}) string {
	return getTypeType(reflect.TypeOf(i))
}

func parseStruct(structType reflect.Type, addParam func(type_ string, name string, desc string)) {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		fmt.Fprintln(os.Stderr, "暂不支持非 struct")
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		type_ := getTypeType(field.Type)

		// 不写 json tag 的话，就用属性类型代替
		name := field.Tag.Get("json")
		if name == "" {
			name = field.Name
		}

		desc := field.Tag.Get("desc")

		addParam(type_, name, desc)
	}
}

func (dg *docGenerator) Bind(i interface{}) error {
	// 传进来的一定是个 struct 指针
	parseStruct(reflect.TypeOf(i).Elem(), docGen.currentAPI().addJSONParam)
	return nil
}

func (dg *docGenerator) JSON(code int, i interface{}) error {
	if code != http.StatusOK {
		// TODO: panic(code)
	}

	data, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		panic(err)
	}
	docGen.currentAPI().returnExampleJSON = string(data)

	switch val := i.(type) {
	case map[string]interface{}:
		if len(val) > 0 {
			for key, value := range val {
				switch val := value.(type) {
				case map[string]interface{}:
					for name, v := range val {
						docGen.currentAPI().addReturnParam(getType(v), name, "")
					}
				default:
					docGen.currentAPI().addReturnParam(getType(val), key, "")
				}
			}
		}
	default:
		// 否则是个 struct 或 struct 指针
		parseStruct(reflect.TypeOf(val), docGen.currentAPI().addReturnParam)
	}

	return nil
}

package example

import (
	"github.com/EndlessCheng/echo-autodoc"
	"github.com/labstack/echo"
	"net/http"
	"time"
	"io/ioutil"
	"strconv"
)

// This is the example.
func runHTTPServer() error {
	e := echo.New()
	setHTTPHandler(autodoc.NewGroup(e)) // 原代码：setHTTPHandler(e)
	return e.Start(":8008")
}

func setHTTPHandler(g autodoc.GroupInterface) { // 原代码：*echo.Echo
	h := &handler{}

	// [skip gen]
	g.GET("/", h.index)

	api := g.Group("/api")

	// 上传一本书
	api.POST("/upload_book", h.uploadBook)

	// 获取一本书的信息
	// 相关出版社有：
	// - 机械工业出版社
	// - 电子工业出版社
	// - 人民邮电出版社
	api.GET("/book", h.getBook)

	// 更新一本书的信息
	// POST 时注意作者需要存到一个数组中
	api.GET("/update_book", h.updateBook)
}

type handler struct {
}

func (h *handler) index(c echo.Context) error {
	return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
}

func (h *handler) uploadBook(c echo.Context) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	f, err := fh.Open()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fileName := c.FormValue("file_name") // 文件名
	if fileName == "" {
		fileName = fh.Filename
	}

	deltaStr := c.FormValue("delta") // int, 偏移量
	_, err = strconv.Atoi(deltaStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &struct {
		Content string `json:"content" desc:"上传的内容"`
	}{string(content)})
}

func (h *handler) getBook(c echo.Context) error {
	c.QueryParam("isbn") // 书的 ISBN
	return c.JSON(http.StatusOK, &exampleBook)
}

func (h *handler) updateBook(c echo.Context) error {
	d := Book{}
	if err := c.Bind(&d); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &struct {
		ErrorCode int    `json:"error_code" desc:"错误码"`
		ErrorMsg  string `json:"error_msg" desc:"错误信息"`
	}{1000, "OK"})
}

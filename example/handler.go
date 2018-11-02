package example

import (
	"github.com/EndlessCheng/echo-autodoc"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

// This is the example.
func runHTTPServer() error {
	e := echo.New()
	setHandlers(autodoc.NewGroup(e)) // 原代码：setHandlers(e)
	return e.Start(":8008")
}

func setHandlers(g autodoc.GroupInterface) { // 原代码：*echo.Echo
	h := &handler{}

	// [skip gen]
	g.GET("/", h.index)

	api := g.Group("/api")

	// 获取一本书的信息
	// 相关出版社有：
	// - 机械工业出版社
	// - 电子工业出版社
	// - 人民邮电出版社
	api.GET("/book", h.getBook)

	// 添加一本书
	// POST 时注意作者需要存到一个数组中
	api.POST("/add_book", h.addBook)
}

type handler struct {
}

func (h *handler) index(c echo.Context) error {
	return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
}

func (h *handler) getBook(c echo.Context) error {
	isbn := c.QueryParam("isbn")
	return c.JSON(http.StatusOK, &Book{ISBN: isbn, Authors: []Author{{}}})
}

func (h *handler) addBook(c echo.Context) error {
	d := Book{}
	if err := c.Bind(&d); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

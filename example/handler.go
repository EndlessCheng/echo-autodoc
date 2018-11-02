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
	// [skip gen]
	g.GET("/", index)

	api := g.Group("/api")

	// 获取一本书的信息
	api.GET("/book", getBook)

	// 添加一本书
	api.POST("/add_book", addBook)
}

func index(c echo.Context) error {
	return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
}

func getBook(c echo.Context) error {
	isbn := c.QueryParam("isbn")
	return c.JSON(http.StatusOK, &Book{ISBN: isbn})
}

func addBook(c echo.Context) error {
	d := Book{}
	if err := c.Bind(&d); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

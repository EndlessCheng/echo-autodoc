package example

type Book struct {
	ISBN        string   `json:"isbn" desc:"ISBN"`
	Name        string   `json:"name" desc:"书名"`
	Authors     []Author `json:"authors" desc:"作者"`
	Press       string   `json:"press" desc:"出版社"`
	PublishDate string   `json:"publish_date" desc:"出版日期"`
	Price       float64  `json:"price" desc:"定价"`
}

type Author struct {
	Name string `json:"name" desc:"姓名"`
	Age  int    `json:"age" desc:"年龄"`
	Sex  int    `json:"sex" desc:"性别（0-男，1-女）"`
}

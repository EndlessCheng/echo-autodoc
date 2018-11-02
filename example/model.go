package example

type Book struct {
	ISBN        string   `json:"isbn" desc:"ISBN"`
	Title       string   `json:"title" desc:"书名"`
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

var exampleBook = Book{
	ISBN:  "9780262033848",
	Title: "Introduction to Algorithms, 3rd Edition",
	Authors: []Author{
		{
			Name: "Thomas H. Cormen",
			Age:  62,
			Sex:  0,
		},
		{
			Name: "Charles E. Leiserson",
			Age:  64,
			Sex:  0,
		},
		{
			Name: "Ronald L. Rivest",
			Age:  71,
			Sex:  0,
		},
		{
			Name: "Clifford Stein",
			Age:  52,
			Sex:  0,
		},
	},
	Press:       "The MIT Press",
	PublishDate: "2009-07-31",
	Price:       94.00,
}

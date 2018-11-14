package example

import (
	"time"
)

type Book struct {
	ID        HashID   `json:"id" desc:"书籍 ID"`
	CreatedAt JSONTime `json:"created_at" desc:"创建时间"`

	ISBN        string   `json:"isbn" desc:"ISBN"`
	Title       string   `json:"title" desc:"书名"`
	Authors     []Author `json:"authors" desc:"作者"`
	Press       string   `json:"press" desc:"出版社"`
	PublishDate string   `json:"publish_date" desc:"出版日期"`
	Price       float64  `json:"price" desc:"定价"`

	Ignored bool `json:"-"`
}

type Author struct {
	ID        HashID    `json:"id" desc:"作者 ID"`
	CreatedAt time.Time `json:"created_at" desc:"创建时间"`

	Name string `json:"name" desc:"姓名"`
	Age  int    `json:"age" desc:"年龄"`
	Sex  int    `json:"sex" desc:"性别（0-男，1-女）"`
}

var exampleBook = Book{
	ID:        1,
	CreatedAt: JSONTime(time.Now()),
	ISBN:      "9780262033848",
	Title:     "Introduction to Algorithms, 3rd Edition",
	Authors: []Author{
		{
			ID:        1,
			CreatedAt: time.Now(),
			Name:      "Thomas H. Cormen",
			Age:       62,
			Sex:       0,
		},
		{
			ID:        2,
			CreatedAt: time.Now(),
			Name:      "Charles E. Leiserson",
			Age:       64,
			Sex:       0,
		},
		{
			ID:        3,
			CreatedAt: time.Now(),
			Name:      "Ronald L. Rivest",
			Age:       71,
			Sex:       0,
		},
		{
			ID:        4,
			CreatedAt: time.Now(),
			Name:      "Clifford Stein",
			Age:       52,
			Sex:       0,
		},
	},
	Press:       "The MIT Press",
	PublishDate: "2009-07-31",
	Price:       94.00,
}

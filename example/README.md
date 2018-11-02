# 接口文档

## HTTP 接口

### 获取一本书的信息

`GET /api/book`

相关出版社有：
- 机械工业出版社
- 电子工业出版社
- 人民邮电出版社

URL 参数

|参数|描述|取值|
|-------|--------|--------|
|isbn|ISBN|string|

返回字段

|参数|描述|取值|
|-------|--------|--------|
|isbn|ISBN|string|
|title|书名|string|
|authors|作者|object array|
|authors.name|姓名|string|
|authors.age|年龄|int|
|authors.sex|性别（0-男，1-女）|int|
|press|出版社|string|
|publish_date|出版日期|string|
|price|定价|float|

返回示例：
```json
{
	"isbn": "9780262033848",
	"title": "Introduction to Algorithms, 3rd Edition",
	"authors": [
		{
			"name": "Thomas H. Cormen",
			"age": 62,
			"sex": 0
		},
		{
			"name": "Charles E. Leiserson",
			"age": 64,
			"sex": 0
		},
		{
			"name": "Ronald L. Rivest",
			"age": 71,
			"sex": 0
		},
		{
			"name": "Clifford Stein",
			"age": 52,
			"sex": 0
		}
	],
	"press": "The MIT Press",
	"publish_date": "2009-07-31",
	"price": 94
}
```


### 添加一本书

`POST /api/add_book`

POST 时注意作者需要存到一个数组中

JSON 参数

|参数|描述|取值|
|-------|--------|--------|
|isbn|ISBN|string|
|title|书名|string|
|authors|作者|object array|
|authors.name|姓名|string|
|authors.age|年龄|int|
|authors.sex|性别（0-男，1-女）|int|
|press|出版社|string|
|publish_date|出版日期|string|
|price|定价|float|

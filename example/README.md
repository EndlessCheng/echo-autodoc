# 接口文档

## HTTP 接口

### 上传一本书

`POST /api/upload_book`

表单参数

|参数|描述|取值|
|-------|--------|--------|
|file|上传的文件|file|
|file_name|文件名|string|
|delta|偏移量|int|

返回字段

|参数|描述|取值|
|-------|--------|--------|
|content|上传的内容|string|

返回示例：
```json
{
	"content": "A"
}
```


### 获取一本书的信息

`GET /api/book`

相关出版社有：
- 机械工业出版社
- 电子工业出版社
- 人民邮电出版社

URL 参数

|参数|描述|取值|
|-------|--------|--------|
|isbn|书的 ISBN|string|

返回字段

|参数|描述|取值|
|-------|--------|--------|
|id|书籍 ID|string|
|created_at|创建时间|string|
|isbn|ISBN|string|
|title|书名|string|
|authors|作者|object array|
|authors.id|作者 ID|string|
|authors.created_at|创建时间|string|
|authors.name|姓名|string|
|authors.age|年龄|int|
|authors.sex|性别（0-男，1-女）|int|
|press|出版社|string|
|publish_date|出版日期|string|
|price|定价|float|

返回示例：
```json
{
	"id": "xQ3Wmvrg",
	"created_at": "2018-11-14 11:31:35",
	"isbn": "9780262033848",
	"title": "Introduction to Algorithms, 3rd Edition",
	"authors": [
		{
			"id": "xQ3Wmvrg",
			"created_at": "2018-11-14T11:31:35.026345+08:00",
			"name": "Thomas H. Cormen",
			"age": 62,
			"sex": 0
		},
		{
			"id": "Rav4N3k2",
			"created_at": "2018-11-14T11:31:35.026345+08:00",
			"name": "Charles E. Leiserson",
			"age": 64,
			"sex": 0
		},
		{
			"id": "8rjZe3xL",
			"created_at": "2018-11-14T11:31:35.026345+08:00",
			"name": "Ronald L. Rivest",
			"age": 71,
			"sex": 0
		},
		{
			"id": "GV67K6zN",
			"created_at": "2018-11-14T11:31:35.026346+08:00",
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


### 更新一本书的信息

`GET /api/update_book`

POST 时注意作者需要存到一个数组中

JSON 参数

|参数|描述|取值|
|-------|--------|--------|
|id|书籍 ID|string|
|created_at|创建时间|string|
|isbn|ISBN|string|
|title|书名|string|
|authors|作者|object array|
|authors.id|作者 ID|string|
|authors.created_at|创建时间|string|
|authors.name|姓名|string|
|authors.age|年龄|int|
|authors.sex|性别（0-男，1-女）|int|
|press|出版社|string|
|publish_date|出版日期|string|
|price|定价|float|

返回字段

|参数|描述|取值|
|-------|--------|--------|
|error_code|错误码|int|
|error_msg|错误信息|string|

返回示例：
```json
{
	"error_code": 1000,
	"error_msg": "OK"
}
```

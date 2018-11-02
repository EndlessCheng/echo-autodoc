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
|name|书名|string|
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
	"isbn": "1",
	"name": "",
	"authors": [
		{
			"name": "",
			"age": 0,
			"sex": 0
		}
	],
	"press": "",
	"publish_date": "",
	"price": 0
}
```


### 添加一本书

`POST /api/add_book`

POST 时注意作者需要存到一个数组中

JSON 参数

|参数|描述|取值|
|-------|--------|--------|
|isbn|ISBN|string|
|name|书名|string|
|authors|作者|object array|
|authors.name|姓名|string|
|authors.age|年龄|int|
|authors.sex|性别（0-男，1-女）|int|
|press|出版社|string|
|publish_date|出版日期|string|
|price|定价|float|

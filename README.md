# gorm-pagination

## Usage

```bash
go get github.com/Hironaga06/gorm-pagination
```

```go
type User struct {
	gorm.Model
	Name string
}

var users []User
db = db.Where("id > ?", 0)

var users []User
p := pagination.New(db, 1, 3, []string{"id asc"}, &users, true)
result, err := p.Paging()
if err != nil {
	panic(err)
}
fmt.Println(result)
```

## With Gin

```go
r := gin.Default()
r.GET("/", func(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "1"))
	if err != nil {
		panic(err)
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "3"))
	if err != nil {
		panic(err)
	}
	var users []User
	p := pagination.New(db, offset, limit, []string{"id asc"}, &users, true)
	result, err := p.Paging()
	if err != nil {
		panic(err)
	}
	c.JSON(200, result)
})
r.Run()
```

## License
[MIT](LICENSE)
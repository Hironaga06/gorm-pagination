package main

import (
	"fmt"
	"strconv"

	pagination "github.com/Hironaga06/gorm-pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

func main() {
	db, err := gorm.Open(postgres.Open("example.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
	var count int64
	if err := db.Model(User{}).Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		users := []User{
			User{Name: "a"},
			User{Name: "b"},
			User{Name: "c"},
			User{Name: "d"},
			User{Name: "e"},
			User{Name: "f"},
			User{Name: "g"},
			User{Name: "h"},
		}
		if err := db.Create(&users).Error; err != nil {
			panic(err)
		}
		fmt.Println("insert complite!")
	}

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
}

package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

type PostRequest struct {
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func main() {
	dsn := "root:coppklmja!BWZ@tcp(127.0.0.1:3306)/items?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移表结构
	type Post struct {
		Content   string    `gorm:"column:content" json:"content"`
		UserID    int       `gorm:"column:user_id" json:"user_id"`
		CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	}

	if err := db.AutoMigrate(&Post{}); err != nil {
		panic(fmt.Sprintf("AutoMigrate failed: %v", err))
	}

	r := gin.Default()

	r.POST("/api/student/post", func(c *gin.Context) {
		var req PostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  "Invalid request body",
			})
			return
		}

		// 验证必填字段
		if req.Content == "" {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  "Content is required",
			})
			return
		}
		if req.UserID <= 0 {
			c.JSON(400, gin.H{
				"code": 400,
				"msg":  "UserID must be a positive number",
			})
			return
		}

		post := Post{
			Content: req.Content,
			UserID:  req.UserID,
		}

		if err := db.Create(&post).Error; err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "Failed to save post",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "success",
		})
	})

	fmt.Println("Server started at http://localhost:8080")
	r.Run(":8080")
}
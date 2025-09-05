package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

type UpdateRequest struct {
	UserID  int    `json:"user_id"`
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type Post struct {
	ID      int    `gorm:"primaryKey"`
	UserID  int    `gorm:"column:user_id"`
	Content string `gorm:"column:content"`
}

func main() {
	r := gin.Default()

	dsn := "root:coppklmja!BWZ@tcp(127.0.0.1:3306)/items?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	r.PUT("/api/student/post", func(c *gin.Context) {
		var req UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, Response{Code: 400, Msg: "Invalid JSON format"})
			return
		}

		// 验证必填字段
		if req.UserID == 0 || req.PostID == 0 || req.Content == "" {
			c.JSON(400, Response{Code: 400, Msg: "Missing required fields"})
			return
		}

		// 查询帖子是否存在且用户权限验证
		var post Post
		if err := db.Where("id = ?", req.PostID).First(&post).Error; err != nil {
			c.JSON(404, Response{Code: 404, Msg: "Post not found"})
			return
		}

		if post.UserID != req.UserID {
			c.JSON(403, Response{Code: 403, Msg: "Unauthorized"})
			return
		}

		// 更新帖子内容
		post.Content = req.Content
		if err := db.Save(&post).Error; err != nil {
			c.JSON(500, Response{Code: 500, Msg: "Failed to update post"})
			return
		}

		c.JSON(200, Response{Code: 200, Data: nil, Msg: "success"})
	})

	r.Run(":8080")
}
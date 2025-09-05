package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
	PostID  int    `json:"post_id"`
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

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database connection")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	r := gin.Default()

	r.DELETE("/api/student/post", func(c *gin.Context) {
		// 解析请求参数
		userID := c.Query("user_id")
		postID := c.Query("post_id")

		// 验证参数
		if userID == "" || postID == "" {
			c.JSON(400, Response{Code: 400, Msg: "Missing user_id or post_id"})
			return
		}

		// 查询帖子是否存在
		var post Post
		if err := db.Where("post_id = ?", postID).First(&post).Error; err != nil {
			c.JSON(404, Response{Code: 404, Msg: "Post not found"})
			return
		}

		// 验证用户权限
		if post.UserID != userID {
			c.JSON(403, Response{Code: 403, Msg: "Unauthorized to delete this post"})
			return
		}

		// 删除帖子
		if err := db.Where("post_id = ?", postID).Delete(&Post{}).Error; err != nil {
			c.JSON(500, Response{Code: 500, Msg: "Failed to delete post"})
			return
		}

		c.JSON(200, Response{Code: 200, Msg: "Post deleted successfully"})
	})

	fmt.Println("Server started at http://localhost:8080")
	r.Run(":8080")
}

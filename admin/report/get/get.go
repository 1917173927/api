package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"time"
)

type Accounts struct {
	UserID   int    `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username string `gorm:"column:username" json:"username"`
	UserType int    `gorm:"column:user_type" json:"user_type"`
}

type Reporteds struct {
	gorm.Model
	UserID  int    `gorm:"column:user_id" json:"user_id"`
	PostID  int    `gorm:"column:post_id" json:"post_id"`
	Content string `gorm:"column:content" json:"content"`
	Reason  string `gorm:"column:reason" json:"reason"`
	Status  int    `gorm:"column:status" json:"status"`
}

type Report struct {
	ReportID int    `json:"report_id"`
	Username string `json:"username"`
	PostID   int    `json:"post_id"`
	Content  string `json:"content"`
	Reason   string `json:"reason"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func main() {
	// 初始化数据库连接
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
	r.GET("/api/admin/report", func(c *gin.Context) {
		// 获取请求参数
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(400, Response{
				Code: 400,
				Msg:  "Missing user_id parameter",
			})
			return
		}

		// 检查用户是否存在
		var account Accounts
		if err := db.Where("user_id = ?", userID).First(&account).Error; err != nil {
			c.JSON(404, Response{
				Code: 404,
				Msg:  "用户不存在",
			})
			return
		}

		// 检查用户权限
		if account.UserType != 2 {
			c.JSON(403, Response{
				Code: 403,
				Msg:  "没有管理员权限",
			})
			return
		}

		// 查询未审批的举报(status=0)
		var reporteds []Reporteds
		if err := db.Where("status = 0").Find(&reporteds).Error; err != nil {
			c.JSON(500, Response{
				Code: 500,
				Msg:  "数据库查询失败",
			})
			return
		}

		// 返回原始记录
		c.JSON(200, Response{
			Code: 0,
			Msg:  "success",
			Data: map[string]interface{}{
				"reports": reporteds,
			},
		})
	})

	fmt.Println("Server started at http://localhost:8080")
	r.Run(":8080")
}
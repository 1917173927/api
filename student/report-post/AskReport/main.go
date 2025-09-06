package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

type Reporteds struct {
	gorm.Model
	UserID  int    `gorm:"column:user_id" json:"user_id"`
	PostID  int    `gorm:"column:post_id" json:"post_id"`
	Content string `gorm:"column:content" json:"content"`
	Reason  string `gorm:"column:reason" json:"reason"`
	Status  int    `gorm:"column:status" json:"status"`
}

type ReportResult struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
	Reason  string `json:"reason"`
	Status  int    `json:"status"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
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
	r.GET("/api/student/report-post", func(c *gin.Context) {
		// 获取请求参数
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(400, gin.H{
				"msg": "Missing user_id parameter",
			})
			return
		}

		// 查询数据库
		var reports []Reporteds
		if err := db.Where("user_id = ?", userID).Find(&reports).Error; err != nil {
			c.JSON(500, Response{
				Code: 500,
				Data: gin.H{},
				Msg:  "数据库查询失败",
			})
			return
		}

		// 转换结果格式
		reportList := make([]ReportResult, 0)
		for _, r := range reports {
			reportList = append(reportList, ReportResult{
				PostID:  r.PostID,
				Content: r.Content,
				Reason:  r.Reason,
				Status:  r.Status,
			})
		}

		// 返回响应
		c.JSON(200, Response{
			Code: 200,
			Data: map[string]interface{}{
				"report_list": reportList,
			},
			Msg: "success",
		})
	})

	fmt.Println("Server started at http://localhost:8080")
	r.Run(":8080")
}
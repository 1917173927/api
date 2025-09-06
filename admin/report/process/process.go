package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"time"
)

type Posts struct {
	PostID  int    `gorm:"primaryKey;column:post_id" json:"post_id"`
	UserID  int    `gorm:"column:user_id" json:"user_id"`
	Content string `gorm:"column:content" json:"content"`
}

type Reporteds struct {
	gorm.Model
	UserID  int    `gorm:"column:user_id" json:"user_id"`
	PostID  int    `gorm:"column:post_id" json:"post_id"`
	Content string `gorm:"column:content" json:"content"`
	Reason  string `gorm:"column:reason" json:"reason"`
	Status  int    `gorm:"column:status" json:"status"`
}

type Accounts struct {
	UserID   int    `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username string `gorm:"column:username" json:"username"`
	UserType int    `gorm:"column:user_type" json:"user_type"`
}

type ApprovalRequest struct {
	UserID   int `json:"user_id"`
	ReportID int `json:"report_id"`
	Approval int `json:"approval"`
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
	r.POST("/api/admin/report", func(c *gin.Context) {
		// 解析请求参数
		var req ApprovalRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, Response{
				Code: 400,
				Msg:  "Invalid request body",
			})
			return
		}

		// 检查用户权限
		var account Accounts
		if err := db.Where("user_id = ? AND user_type = 2", req.UserID).First(&account).Error; err != nil {
			c.JSON(403, Response{
				Code: 403,
				Msg:  "没有管理员权限",
			})
			return
		}

		// 处理审批
		var msg string
		var code int
		if req.Approval == 1 {
			// 获取举报记录
			var reported Reporteds
			if err := db.First(&reported, req.ReportID).Error; err != nil {
				c.JSON(404, Response{
					Code: 404,
					Msg:  "举报记录不存在",
				})
				return
			}

			// 删除帖子
			if err := db.Where("post_id = ?", reported.PostID).Delete(&Posts{}).Error; err != nil {
				c.JSON(500, Response{
					Code: 500,
					Msg:  "删除帖子失败",
				})
				return
			}

			// 删除举报记录
			if err := db.Delete(&reported).Error; err != nil {
				c.JSON(500, Response{
					Code: 500,
					Msg:  "删除举报记录失败",
				})
				return
			}

			msg = "帖子已删除"
			code = 1
		} else {
			// 更新举报状态为已处理
			if err := db.Model(&Reporteds{}).Where("id = ?", req.ReportID).Update("status", 1).Error; err != nil {
				c.JSON(500, Response{
					Code: 500,
					Msg:  "更新举报状态失败",
				})
				return
			}

			msg = "举报已拒绝"
			code = 2
		}

		// 返回响应
		c.JSON(200, Response{
			Code: code,
			Msg:  msg,
			Data: nil,
		})
	})

	fmt.Println("Server started at http://localhost:8080")
	r.Run(":8080")
}
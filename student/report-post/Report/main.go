package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"time"
	"fmt"
)

type Reported struct {
	gorm.Model
	UserID int    `gorm:"column:user_id" json:"user_id"`
	PostID int    `gorm:"column:post_id" json:"post_id"`
	Reason string `gorm:"column:reason" json:"reason"`
}

func main() {
	r := gin.Default()

	r.POST("/api/student/report-post", func(c *gin.Context) {
		var data Reported

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{
				"code": 400,
				"data": nil,
				"msg":  "绑定数据失败",
			})
			return
		}
		dsn := "root:coppklmja!BWZ@tcp(127.0.0.1:3306)/items?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		fmt.Println(err)
		fmt.Println(db)
		//设置连接池
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get database connection")
		}	
		// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
		sqlDB.SetMaxIdleConns(10)
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(100)
		// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
		sqlDB.SetConnMaxLifetime(10 * time.Second) //10s
		if err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"data": nil,
				"msg":  "数据库连接失败",
			})
			return
		}

		if err := db.AutoMigrate(&Reported{}); err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"data": nil,
				"msg":  "数据库迁移失败",
			})
			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"data": nil,
				"msg":  "数据插入失败",
			})
			return
		}

		c.JSON(200, gin.H{
			"code": 200,
			"data": nil,
			"msg":  "success",
		})
	})

	r.Run(":8080")
}
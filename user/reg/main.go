package main

import (
	"fmt"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type accounts struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserType int    `json:"user_type"` // 1: 学生, 2: 管理员
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func main() {

	//数据库连接
	dsn := "root:coppklmja!BWZ@tcp(127.0.0.1:3306)/items?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
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
	//Gin初始化
	r := gin.Default()

	// 注册路由
	r.POST("/api/user/reg", func(c *gin.Context) {
		var req accounts
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Response{Code: 400, Msg: "Invalid request body"})
			return
		}

		// 验证必填字段
		if req.Username == "" || req.Password == "" || (req.UserType != 1 && req.UserType != 2) {
			c.JSON(http.StatusBadRequest, Response{Code: 400, Msg: "Missing or invalid fields"})
			return
		}

		// 数据库操作
		if err := db.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{Code: 500, Msg: "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, Response{Code: 200, Data: req, Msg: "User registered successfully"})
	})

	//端口号
	PORT := "8080"
	fmt.Println("Server started at http://localhost:" + PORT)
	r.Run(":" + PORT)
}
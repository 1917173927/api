//MySQL+GORM+GIN大学习
//curd
package main
import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
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
 
	//结构体
	type List struct {
		gorm.Model
		Name string `gorm:"type:varchar(100);not null" json:"name"`
		Age  int    `gorm:"column:age" json:"age"`
	}

	// db.AutoMigrate(&List{})
	// if err := db.AutoMigrate(&List{}); err != nil {
    // fmt.Printf("AutoMigrate failed: %v\n", err)
	// }

	//接口
	r := gin.Default()

	// //TEST
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Success",
	// 	})
	// })
	//增
	r.POST("/user/add",func(c *gin.Context) {
		var data List

		err := c.ShouldBindJSON(&data)

		//判断绑定是否错误
		if err != nil {
			c.JSON(200, gin.H{
				"msg": "添加失败",
				"data": gin.H{},
				"code": 400,
			})
			
		}else{

			//数据库操作
			db.Create(&data)//创建一条数据

			c.JSON(200, gin.H{
				"msg": "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})
	//删
	r.DELETE("/user/delete/:id",func(c *gin.Context) {
		var data []List
		//接受id
		id := c.Param("id")
		//判断id是否存在
		db.Where("id=?",id).Find(&data)
		//id存在删除，不存在报错
		if len(data)==0{
			c.JSON(200,gin.H{
				"msg":"IdNotFound",
				"code":400,
			})
		}else {
			//删除
			db.Where("id=?",id).Delete(&data)
			c.JSON(200,gin.H{
				"msg":"删除成功",
				"code":200,
			})
		}
	})
	//改
	r.PUT("/user/update/:id",func(c *gin.Context) {
		var data List
		//接受id
		id:=c.Param("id")
		//判断id是否存在
		db.Select("id").Where("id=?",id).Find(&data)
		if data.ID==0{
			c.JSON(200,gin.H{
				"msg":"IdNotFound",
				"code":400,
			})
		}else {
			//接受参数
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(200,gin.H{
					"msg":"修改错误",
					"code":400,
				})
			}else{
				db.Where("id=?",id).Updates(&data)
				c.JSON(200,gin.H{
					"msg":"修改成功",
					"code":200,
				})
			}
		}
	})
	//查
	r.GET("/user/list/:name",func(c *gin.Context) {
		name := c.Param("name")
		var dataList []List
		db.Where("name=?",name).Find(&dataList)
		if len(dataList)==0{
			c.JSON(200,gin.H{
				"msg":"NotFound",
				"code":400,
				"data":gin.H{},
			})
		}else {
			c.JSON(200,gin.H{
				"msg":"查询成功",
				"code":200,
				"data":dataList,
			})
		}
	})

	//端口号
	PORT := "8080"
	r.Run(":" + PORT)

}
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code    int         `json:"code"`
	Data    UserData    `json:"data"`
	Msg     string      `json:"msg"`
}

type UserData struct {
	UserID   int `json:"user_id"`
	UserType int `json:"user_type"`
}

type accounts struct {
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	UserID   int    `gorm:"column:user_id"`
	UserType int    `gorm:"column:user_type"`
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
  dsn := "root:coppklmja!BWZ@tcp(127.0.0.1:3306)/items?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  fmt.Println(db, err)
	http.HandleFunc("/api/user/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 查询数据库中的用户记录
		var account accounts
		if err := db.Where("username = ?", req.Username).First(&account).Error; err != nil {
			resp := LoginResponse{
				Code: 400,
				Data: UserData{},
				Msg:  "Reject: User not found",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		// 验证密码
		if account.Password != req.Password {
			resp := LoginResponse{
				Code: 400,
				Data: UserData{},
				Msg:  "Reject: Invalid password",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		// 验证通过
		resp := LoginResponse{
			Code: 200,
			Data: UserData{
				UserID:   account.UserID,
				UserType: account.UserType,
			},
			Msg: "Success",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
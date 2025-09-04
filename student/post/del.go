package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func main() {
	http.HandleFunc("/api/student/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 解析请求参数
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
		if err != nil {
			http.Error(w, "Invalid post_id", http.StatusBadRequest)
			return
		}

		// 模拟验证用户权限（实际项目中需查询数据库）
		if userID != 1 { // 假设用户ID为1有权限
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 模拟删除帖子（实际项目中需操作数据库）
		fmt.Printf("Post %d deleted by user %d\n", postID, userID)

		// 返回响应
		resp := Response{
			Code: 200,
			Data: nil,
			Msg:  "success",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
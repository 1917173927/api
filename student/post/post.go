package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PostRequest struct {
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func main() {
	http.HandleFunc("/api/student/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req PostRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 验证必填字段
		if req.Content == "" || req.UserID == 0 {
			http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
			return
		}

		// 记录发帖时间（可替换为实际数据库操作）
		postTime := time.Now()
		fmt.Printf("User %d posted at %v\n", req.UserID, postTime)

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
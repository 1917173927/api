package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReportRequest struct {
	UserID int    `json:"user_id"`
	PostID int    `json:"post_id"`
	Reason string `json:"reason"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func main() {
	http.HandleFunc("/api/student/report", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 读取请求体
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 解析请求体
		var req ReportRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// 验证必填字段
		if req.UserID == 0 || req.PostID == 0 || req.Reason == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// 模拟举报逻辑（实际项目中需记录到数据库）
		fmt.Printf("Report received: User %d reported Post %d for reason: %s\n", req.UserID, req.PostID, req.Reason)

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
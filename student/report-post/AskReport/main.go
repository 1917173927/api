package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	http.HandleFunc("/api/student/report/result", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 获取请求参数
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
			return
		}

		// 模拟查询举报结果（实际项目中需查询数据库）
		reportList := []ReportResult{
			{
				PostID:  2,
				Content: "1233",
				Reason:  "123",
				Status:  0,
			},
		}

		// 返回响应
		resp := Response{
			Code: 200,
			Data: map[string]interface{}{
				"report_list": reportList,
			},
			Msg: "success",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
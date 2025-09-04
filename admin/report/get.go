package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Report struct {
	ReportID int    `json:"report_id"`
	Username string `json:"username"`
	PostID   int    `json:"post_id"`
	Content  string `json:"content"`
	Reason   string `json:"reason"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func main() {
	http.HandleFunc("/api/admin/report", func(w http.ResponseWriter, r *http.Request) {
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

		// 模拟查询未审批的举报帖子列表（实际项目中需查询数据库）
		reportList := []Report{}
		if userID == "1" {
			reportList = []Report{
				{
					ReportID: 1,
					Username: "student1",
					PostID:   101,
					Content:  "违规内容",
					Reason:   "搞黄色",
				},
			}
		}

		// 返回响应
		resp := Response{
			Code: 0,
			Msg:  "success",
			Data: map[string]interface{}{
				"report_list": reportList,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
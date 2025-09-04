package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	http.HandleFunc("/api/admin/report", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 解析请求参数
		var req ApprovalRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 审批逻辑完全由API提交的approval参数控制
		var msg string
		var code int
		if req.Approval == 1 {
			// 审批通过，后续逻辑由API调用者实现
			msg = "Accepted"
			code = 1
		} else {
			// 审批拒绝
			msg = "Rejected"
			code = 2
		}

		// 返回响应
		resp := Response{
			Code: code,
			Msg:  msg,
			Data: nil,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
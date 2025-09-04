package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Post struct {
	ID      int       `json:"id"`
	Content string    `json:"content"`
	UserID  int       `json:"user_id"`
	Time    time.Time `json:"time"`
	Likes   int       `json:"likes"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func main() {
	http.HandleFunc("/api/student/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 模拟帖子数据
		postList := []Post{
			{
				ID:      1,
				Content: "123",
				UserID:  2,
				Time:    time.Now().Add(-24 * time.Hour),
				Likes:   10,
			},
			{
				ID:      2,
				Content: "1233",
				UserID:  3,
				Time:    time.Now().Add(-12 * time.Hour),
				Likes:   10,
			},
		}

		// 返回响应
		resp := Response{
			Code: 200,
			Data: map[string]interface{}{"post_list": postList},
			Msg:  "success",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
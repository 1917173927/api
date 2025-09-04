package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code int      `json:"code"`
	Data UserData `json:"data"`
	Msg  string   `json:"msg"`
}

type UserData struct {
	UserID   int `json:"user_id"`
	UserType int `json:"user_type"`
}

func main() {
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

		var userType int
		if req.Username == "admin" {
			userType = 2 // Admin
		} else {
			userType = 1 // Student
		}

		resp := LoginResponse{
			Code: 200,
			Data: UserData{
				UserID:   1,
				UserType: userType,
			},
			Msg: "success",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

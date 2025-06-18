package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// 定义 POST 请求体结构
type Message struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, this is a GET request!")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "JSON 解析失败", http.StatusBadRequest)
		return
	}

	log.Printf("收到消息: %+v\n", msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"echo":   fmt.Sprintf("收到 %s: %s", msg.Name, msg.Text),
	})
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/post", postHandler)

	port := "8080"
	log.Printf("服务器启动在 http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

func main() {
	// 显示使用多少核
	fmt.Printf("GOMAXPROCS = %d\n", runtime.GOMAXPROCS(0))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 模拟 CPU 密集负载（可调）
		for i := 0; i < 1e5; i++ {
			_ = i * i
		}

		duration := time.Since(start)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong in %v\n", duration)
	})

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}

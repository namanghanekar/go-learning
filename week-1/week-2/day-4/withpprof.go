package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

// 🔴 Memory heavy function
func memoryHeavy() {
	for i := 0; i < 100000; i++ {
		data := make([]byte, 1024*50) // 50KB allocation
		_ = data
	}
}

func main() {

	// ✅ Start pprof server
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	// ✅ Run workload continuously
	for {
		memoryHeavy()
		time.Sleep(1 * time.Second)
	}
}

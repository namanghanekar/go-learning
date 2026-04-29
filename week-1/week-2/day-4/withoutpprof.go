package main

import (
	"fmt"
	"runtime"
	"time"
)

// 🔴 This function creates heavy memory allocation
func memoryHeavyWithoutPprof() {
	for i := 0; i < 100000; i++ {
		data := make([]byte, 1024*50) // 50KB allocation
		_ = data
	}
}

// 🟢 This function prints memory stats
func printMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v KB\n", m.Alloc/1024)
	fmt.Printf("TotalAlloc = %v KB\n", m.TotalAlloc/1024)
	fmt.Printf("Sys = %v KB\n", m.Sys/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
	fmt.Println("------")
}

func main() {

	// ✅ Run continuously
	for {
		memoryHeavyWithoutPprof() // Step 1: allocate memory

		printMem() // Step 2: check memory AFTER allocation

		time.Sleep(1 * time.Second) // Step 3: observe changes
	}
}

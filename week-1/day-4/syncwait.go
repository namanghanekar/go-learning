package main

import (
	"fmt"
	"sync"
)

func task(wg *sync.WaitGroup, id int) {
	defer wg.Done() //execution of the task is done then only it will call the Done() method
	fmt.Println("Task", id, "done")
}
func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go task(&wg, i)
	}
	wg.Wait()
	fmt.Println("All tasks completed")
}

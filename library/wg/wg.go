package wg

import (
	"sync"
)

var (
	wg = new(sync.WaitGroup)

	// true 的时候才可以添加任务，关掉后不可以
	status bool = true
)

func Add(i int) bool {
	if status == true {
		wg.Add(i)
		return true
	}
	return false
}

func Wait() {
	wg.Wait()
}

func Done() {
	wg.Done()
}

func Close() {
	status = false
}

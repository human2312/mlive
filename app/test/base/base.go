package main

/**
 * @Author: Lemyhello
 * @Description:
 * @File:  base
 * @Version: X.X.X
 * @Date: 2020/3/25 下午2:07
 */

import (
	"fmt"
	"mlive/app/test"
	"mlive/library/graceful"
	"mlive/library/wg"
	"time"
)

func main() {
	test.ConfigInit()
	wg.Add(1)
	go router()
	//监听信号并处理
	go func() {
		test.SignalDeal()
	}()

	// 阻塞
	wg.Wait()
}

//router 消费进程
func router()  {
	fmt.Println(222)
	for {
		//防止堵塞
		go graceful.Put(false)

		time.Sleep(2 * time.Second)
		fmt.Println(111)
		fmt.Println(graceful.Pointer())
		isKill := graceful.Get() //停止信号传递到管道
		if isKill == true {
			wg.Done()
			break
		}
	}
	//time.Sleep(10 * time.Second)
	fmt.Println("任务执行结束")
	return
}

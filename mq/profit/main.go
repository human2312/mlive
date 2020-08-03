package main

import (
	"fmt"
	"mlive/library/wg"
	"mlive/mq"
	"mlive/service/profit"
	"mlive/util"
	"runtime/debug"
)

func main() {
	mq.ConfigInit()
	wg.Add(1)
	go router()
	//监听信号并处理
	go func() {
		mq.SignalDeal()
	}()
	// 阻塞
	wg.Wait()
}

//router 消费进程
func router() {
	//所有panic接收
	defer func() {
		if err := recover(); err != nil {
			util.SendDingDing("mq/profit进程",err,"【堆栈信息】", string(debug.Stack()))
			fmt.Println("mq/profit进程",err,"【堆栈信息】", string(debug.Stack()))
			panic(recover())
		}
	}()

	profit.OrderReceive()
}

package main

import (
	"fmt"
	"mlive/app"
	"mlive/library/wg"
	"mlive/mq"
	"mlive/util"
	"runtime/debug"
)

/**
 * @Author: Lemyhello
 * @Description:
 * @File:  main
 * @Version: X.X.X
 * @Date: 2020/3/25 下午10:33
 */

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
func router()  {
	//所有panic接收
	defer func() {
		if err := recover(); err != nil {
			util.SendDingDing("mq/mqgraceful进程",err,"【堆栈信息】", string(debug.Stack()))
			fmt.Println("mq/mqgraceful进程",err,"【堆栈信息】", string(debug.Stack()))
			panic(recover())
		}
	}()

	app.ReceiveTest()
}
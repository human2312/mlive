package main

/**
 * @Author: ray
 * @Description: 用户信息监听
 * @Date: 2020/3/16 下午03:00
 */

import (
	"fmt"
	"mlive/app"
	"mlive/library/wg"
	"mlive/mq"
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
func router()  {
	//所有panic接收
	defer func() {
		if err := recover(); err != nil {
			util.SendDingDing("mq/userinfo进程",err,"【堆栈信息】", string(debug.Stack()))
			fmt.Println("mq/userinfo进程",err,"【堆栈信息】", string(debug.Stack()))
			panic(recover())
		}
	}()

	app.MqUserInfoReceive()
}
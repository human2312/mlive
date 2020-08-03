package test

import (
	"fmt"
	"mlive/dao"
	"mlive/library/config"
	"mlive/library/graceful"
	"mlive/library/logger"
	"mlive/library/wg"
	"mlive/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ConfigInit() {
	// 初始化配置
	config.Init()
	logger.Init()
	dao.NewDao()
	//mq连接池
	util.NewMQ()
}


// sighup 1
// sigint 2 | ctrl+c
// sigterm 15 | kill pid
// SIGQUIT	3		来自键盘的离开信号
// 9 sigkill 无法被监听，会被强杀 不可以被捕获或忽略的终止信号
// SIGSTOP 19 不能被捕获或忽略的停止信号
func SignalDeal() {
	q := make(chan os.Signal)
	signal.Notify(q)
	for {
		switch <-q {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			go func() {
				//防止进程假死，中断指令一个周期后自动断开阻塞
				time.Sleep( 10 * time.Second)
				fmt.Println("强制进程结束")
				wg.Done()
			}()
			fmt.Println("接收结束指令")
			fmt.Println(graceful.Pointer())
			graceful.Put(true)
			wg.Close()
			return
		default:
			logger.Iprintln("received a signal but ignore")
		}
	}
}

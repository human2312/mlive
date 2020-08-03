package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"mlive/library/config"
	"mlive/library/logger"
	"mlive/library/wg"
	"mlive/router"
)

var (
	svr *http.Server
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化日志
	logger.Init()

	// 设置mode
	gin.SetMode(viper.GetString("gin.mode"))

	// router初始化
	r := router.Init()

	// 初始化server
	svr = &http.Server{
		Addr:         viper.GetString("server.addr"),
		Handler:      r,
		ReadTimeout:  viper.GetDuration("server.readTimeout") * time.Second,
		WriteTimeout: viper.GetDuration("server.writeTimeout") * time.Second,
	}

	// 一个goroutine 跑http
	wg.Add(2)
	go func() {
		svr.ListenAndServe()
		wg.Done()
	}()

	// 一个goroutine 监听信号并处理
	go func() {
		signalDeal()
		wg.Done()
	}()

	// 阻塞
	wg.Wait()

}

// sighup 1
// sigint 2 | ctrl+c
// sigterm 15 | kill pid
// SIGQUIT	3		来自键盘的离开信号
// 9 sigkill 无法被监听，会被强杀 不可以被捕获或忽略的终止信号
// SIGSTOP 19 不能被捕获或忽略的停止信号
func signalDeal() {
	q := make(chan os.Signal)
	signal.Notify(q)
	for {
		switch <-q {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			wg.Close()
			ctx, cancel := context.WithTimeout(context.Background(), (viper.GetDuration("server.readTimeout")+viper.GetDuration("server.writeTimeout"))*time.Second)
			defer cancel()
			if err := svr.Shutdown(ctx); err != nil {
				logger.Iprintln("Server Shutdown:", err)
			} else {
				logger.Iprintln("server shutdown gracefully")
			}
			return
		default:
			logger.Iprintln("received a signal but ignore")
		}
	}
}

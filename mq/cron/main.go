package main

import (
	"fmt"
	"github.com/robfig/cron"
	"mlive/mq"
	cloud "mlive/service/cloud/dao"
	"mlive/util"
	"runtime/debug"
)


// 定时任务 (https://www.cnblogs.com/liuzhongchao/p/9521897.html)

// cron举例说明
// 每隔5秒执行一次：*/5 * * * * ?
// 每隔1分钟执行一次：0 */1 * * * ?
// 每天23点执行一次：0 0 23 * * ?
// 每天凌晨1点执行一次：0 0 1 * * ?
// 每月1号凌晨1点执行一次：0 0 1 1 * ?
// 在26分、29分、33分执行一次：0 26,29,33 * * * ?
// 每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?

func main()  {
	mq.ConfigInit()
	c := cron.New()
	CronTimingCloudInfo(c)
	//启动计划任务
	c.Start()
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()
	select{}
}

// 定时处理云仓-订单
func CronTimingCloudInfo(c *cron.Cron)  {

	//所有panic接收
	defer func() {
		if err := recover(); err != nil {
			util.SendDingDing("mq/cron进程",err,"【堆栈信息】", string(debug.Stack()))
			fmt.Println("mq/cron进程",err,"【堆栈信息】", string(debug.Stack()))
			panic(recover())
		}
	}()

	//AddFunc
	spec := "0 */1 * * * ?"
	//c.AddFunc(spec, func() {
	//	log.Println("cron running:")
	//})
	//AddJob方法
	c.AddJob(spec, cloud.TimingCloudInfo{})
}


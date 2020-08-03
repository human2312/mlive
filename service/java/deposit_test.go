package java

import (
	"mlive/dao"
	"mlive/library/config"
	"mlive/library/logger"
	"testing"
)

func init()  {

	// 初始化配置
	config.Init()
	// 初始化日志
	logger.Init()
	// 链接数据库
	dao.NewDao()
}

func TestSubmitOrderDeposit(t *testing.T)  {

	var userId int64 = 65
	var ouponSn int64 = 111
	var orderType int64 = 1

	res,err := SubmitOrderDeposit(ouponSn,orderType,userId,21,"")
	t.Log("res:",res)
	t.Log("err:",err)
}
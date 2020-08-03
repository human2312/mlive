package php

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

func TestGetInvestorsList(t *testing.T)  {

	info,err := GetInvestorsList()
	t.Log("info:",info.Data.Items[0])
	t.Log("err:",err)
}
package user

import (
	"mlive/dao"
	"mlive/library/config"
	"mlive/library/logger"
	"testing"
)

func initC()  {
	// 初始化配置
	config.Init()
	// 初始化日志
	logger.Init()
	//  db 初始化
	dao.NewDao()

}

func TestUser(t *testing.T)  {

	initC()
	User()
	t.Log("111")
}


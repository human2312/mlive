package dao

import (
	"mlive/dao"
	"mlive/library/config"
	"testing"
	"mlive/library/logger"
	"mlive/util"

)


func init()  {

	// 初始化配置
	config.Init()
	// 链接数据库
	dao.NewDao()
	// 初始化日志
	logger.Init()
	// 初始化mq
	util.NewMQ()
}

func TestSaveLevel(t *testing.T)  {
	//var(
	//	userId int64  = 1
	//	ty     int64  = 1
	//	level  int64  = 1
	//)
	// 初始化日志
	logger.Init()
	//res,err := u.SaveLevel(userId,0,ty,level,0,"","")
	var u User

	res,err := u.SaveLevel(11428, 3, 2, 3, 0, "", "2020040200227425")
	t.Log("res:",res)
	t.Log("err:",err)



}
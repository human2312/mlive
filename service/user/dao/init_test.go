package dao


import (
	"log"
	"mlive/dao"
	"mlive/library/config"
	"mlive/util"
	"testing"
	"mlive/library/logger"
)


func init()  {

	// 初始化配置
	config.Init()
	// 初始化日志
	logger.Init()
	// 链接数据库
	dao.NewDao()
	// 初始化mq
	util.NewMQ()

}

// 修复邀请上级=0的用户
func TestInitInviteData(t *testing.T)  {


	list,err := u.InitInviteData()
	log.Println("list",list)
	log.Println("err",err)
}

// 修复和java赚播数据不一致的用户数据
func TestInitSySnLevel(t *testing.T)  {

	list,err := u.InitSySnLevel()
	log.Println("list",list)
	log.Println("err",err)
}

//修复用户状态
func TestInitSyncStatus(t *testing.T)  {

	list,err := u.InitSyncStatus()
	log.Println("list",list)
	log.Println("err",err)
}

func TestInitSyncMq(t *testing.T)  {

	list,err := u.InitSyncMq()
	log.Println("list",list)
	log.Println("err",err)
}
package dao

import (
	"mlive/dao"
	"mlive/library/config"
	"testing"
)

func init()  {

	// 初始化配置
	config.Init()
	// 链接数据库
	dao.NewDao()
}


func TestSaveReduceCloudStoragNums(t *testing.T)  {

	var(
		userId  int64 = -1
	)

	var c CloudStorage
	info,err :=c.SaveReduceCloudStoragNums(userId,1,1,0,"test",2,0)
	t.Log("info:",info)
	t.Log("err:",err)
}
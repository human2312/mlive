package clear

import (
	"mlive/dao"
	"mlive/library/config"
	"testing"
)

var cdb *ClearDb

func init() {

	// 初始化配置
	config.Init()
	// 链接数据库
	dao.NewDao()
}


func TestToClear(t *testing.T)  {

	status,err := cdb.ToClear()
	t.Log("status:",status)
	t.Log("err:",err)
}
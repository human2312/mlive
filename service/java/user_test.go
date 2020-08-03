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

func TestInfoById(t *testing.T)  {

	var userId int64 = 1
	info,err  := InfoById(userId)
	t.Log("info:",info)
	t.Log("err:",err)
}
//
//func TestGetOrderBuyTypeCount(t *testing.T)  {
//
//	var userId int64 = 65
//	count,err := GetOrderBuyTypeCount(userId)
//	t.Log("count:",count)
//	t.Log("err:",err)
//}

func TestGetAdminUserAccountInfoStatistics(t *testing.T)  {

	var userId int64 = 110
	var AdminToken string = "55gpahio4zcfowwc8fdt5uiv7p9cjlso"
	info,err := GetAdminUserAccountInfoStatistics(userId,AdminToken)
	t.Log("info:",info)
	t.Log("err:",err)
}
//
//func TestGetDepositInfo(t *testing.T)  {
//
//	var depositNo = "2020032700226020"
//	info,err := GetDepositInfo(depositNo)
//	t.Log("info:",info.Data)
//	t.Log("err:",err)
//}

func TestGetXFYLUserInfo(t *testing.T)  {

	var (
		areaCode = "86"
		mobile = "13717276544"
	)
	xfylUserInfo,err :=GetXFYLUserInfo(areaCode,mobile)
	if err != nil || xfylUserInfo.Data.UserId <=0  {
		t.Log("xfylUserInfo:",xfylUserInfo.Data.UserId)
	}
	t.Log("xfylUserInfo:11",xfylUserInfo)
	t.Log("err:",err)
}
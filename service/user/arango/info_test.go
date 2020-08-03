package arango

import (
	"mlive/dao"
	"mlive/library/config"
	"testing"
	"time"
	"mlive/library/logger"
)

func init()  {

	// 初始化配置
	config.Init()
	// 链接数据库
	dao.NewDao()
	// 初始化日志
	logger.Init()
}

var da ArangoUser

func TestCreate(t *testing.T)  {

	tTime		:= time.Now()
	nowTime 	:=  tTime.Format("2006-01-02 15:04:05")

	var (
		Id int64 = 1
		Level int64 = 1
		InviteUpId  int64 = 0
	)
	mapData := map[string]interface{}{
		"id":Id,
		"user_name":"测试",
		"level":Level,
		"invite_up_id":InviteUpId,
		"invite_time":nowTime,
		"create_time":nowTime,
		"update_time":nowTime,
	}
	res ,err := da.Create(mapData)
	t.Log("res:",res)
	t.Log("err:",err)

}

func TestGetTeamList(t *testing.T)  {

	var(
		userId int64 = 11209
		level  int64	= 5
	)

	info,err := da.GetTeamList(userId,level)
	t.Log("info:",info)
	t.Log("err:",err)
}

func TestGetSameLevelReward(t *testing.T)  {

	var(
		userId int64 = 11214
	)
	info,err := da.GetSameLevelReward(userId)
	t.Log("info:",info)
	t.Log("err:",err)
}
func TestGetInviteSameLevel(t *testing.T)  {
	var(
		userId int64 = 11464
		iLevel int64 = 5
	)
	info,err := da.GetInviteSameLevel(userId,iLevel)
	t.Log("info:",info)
	t.Log("err:",err)
}

func TestCreateIndex(t *testing.T)  {

	da.CreateIndex()
}

func TestGetLeapFrogInfo(t *testing.T)  {

	var(
		userId int64 = 410354
	)
	info,err := da.GetLeapFrogInfo(userId)
	t.Log("info:",info)
	t.Log("err:",err)
}
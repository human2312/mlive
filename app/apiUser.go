package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mlive/mq"
	"mlive/service/java"
	"mlive/service/user/arango"
	"mlive/util"
	"strconv"
	"sync"
)


// 获取用户统计信息
func GetUserStatisticsInfo(c *gin.Context)  {

	checkParams   := []interface{}{"userId"}
	params,err 	  := util.GetRawData(c.Request,checkParams)
	if err != nil {
		util.Fail(c, 11111, fmt.Sprintf("%s",err))
		return
	}
	userParams := util.ChectIntFloat(params["userId"])
	if !userParams {
		util.Fail(c, 11111, "parameter fail")
		return
	}
	userId := int64(params["userId"].(float64))

	if userId <= 0 {
		c.JSON(200, gin.H{
			"code": 11111,
			"msg":  "缺少用户id",
			"data": []string{},
		})
		return
	}
	info,_ := user.GetUserInfo(userId)
	if info.Id <= 0{
		c.JSON(200, gin.H{
			"code": 80000,
			"msg":  "用户不存在",
			"data": []string{},
		})
		return
	}
	var (
		totalNum int64 = 0
		list []arango.AgoTeamList = []arango.AgoTeamList{}
		addNum int = 0
	)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		totalNum,_ 	  = ago.GetTeamNum(userId,info.Level)
		wg.Done()
	}(wg)

	go func(wg *sync.WaitGroup) {
		list,_ 	  = ago.GetTeamList(userId,info.Level)
		addNum 	  = len(list)
		wg.Done()
	}(wg)

	wg.Wait()
	if list == nil {
		list = []arango.AgoTeamList{}
	}

	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "",
		"data": map[string]interface{}{
			"teamList": list,
			"totalNum": totalNum,
			"addNum": addNum,
		},
	})


}


type UserTeamList struct {
	UserId 		int64				`json:"userId"`
	Nickname 	string				`json:"nickname"`
	Level	 	int64				`json:"level"`
	InviteUpId  int64				`json:"inviteUpId"`
	InviteCode	string				`json:"inviteCode"`
	TeamNum     int64				`json:"teamNum"`
	InviteNum 	int64				`json:"inviteNum"`
	Status	    int64				`json:"status"`
	CreateTime	string				`json:"createTime"`
	UpdateTime	string				`json:"updateTime"`
	ChildList 	[]*UserTeamList		`json:"children"`
}

// 给联创做自己的邀请链条数状
func GetInviteTree(c *gin.Context)  {

	userToken := c.Request.Header.Get("X-MPMall-Token")
	userId ,err:= java.Token2Id(userToken)
	if userId <= 0 {
		util.Fail(c,80000," user token fail:"+err.Msg)
		return
	}
	var toUserId int64 = 0
	checkParams   	  := []interface{}{"inviteUpId"}
	params,err1 	 	  := util.GetRawData(c.Request,checkParams)
	if err1 != nil {
		util.Fail(c,11111, fmt.Sprintf("%s",err1))
		return
	}
	inviteUpIdParams := util.ChectIntFloat(params["inviteUpId"])
	if !inviteUpIdParams {
		util.Fail(c,11111, "parameter fail")
		return
	}

	userInfo,_ := user.GetUserInfo(userId)
	if userInfo.Id <= 0 {
		util.Fail(c,80000," get user  fail")
		return
	}
	inviteUpId := int64(params["inviteUpId"].(float64))

	if inviteUpId <= 0 {
		toUserId = userInfo.Id
	}

	if userInfo.Level != 5 {
		util.Fail(c,80000," no lian  level ,fail")
		return
	}

	if inviteUpId >0 && userInfo.Id != inviteUpId {
		isSub,_ := ago.CheckIsInviteSubordinate(userInfo.Id,inviteUpId)
		if !isSub {
			util.Fail(c,80000," select inviteUpId tree error")
			return
		}
	}
	var list []arango.MyUserList
	var teamList []UserTeamList

	list,_ = ago.GetInviteList(toUserId,inviteUpId)
	if  len(list) >= 1 {
		for _,v := range list {
			teamList = append(teamList,UserTeamList{
				UserId:v.Id,
				Nickname:v.Nickname,
				Level:v.Level,
				InviteUpId:v.InviteUpId,
				InviteCode:v.InviteCode,
				TeamNum:0,
				InviteNum:0,
				Status:v.Status,
				ChildList: []*UserTeamList{},
				CreateTime:v.CreateTime,
				UpdateTime:v.UpdateTime,
			})
		}
	}
	if len(teamList) > 0 {
		wg := &sync.WaitGroup{}
		for i := 0; i < len(teamList); i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				// 团队人数
				teamList[i].TeamNum, _ = ago.GetTeamNum(teamList[i].UserId, teamList[i].Level)
				//// 直属(邀请人数)人数
				teamList[i].InviteNum, _ = user.GetUserCount(0, teamList[i].UserId)
				wg.Done()
			}(wg,i)
		}
		wg.Wait()
	}else{
		teamList = []UserTeamList{}
	}
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data":map[string]interface{}{
			"items":teamList,
			"total":len(teamList),
		},
	})

}



type MqGeneratedUser struct {
	UserId 				int64 		 `json:"userId"`
}

var (
	mqUserLevelUser          = new(util.RabbitMQ)
)


//发送一条消息
func MqUser(c *gin.Context)  {
	mq.ConfigInit()
	sendJson := MqGeneratedUser{}

	userId,_ := strconv.ParseInt(c.DefaultQuery("userId","0"),10,0)
	if userId <= 0 {
		c.JSON(200, gin.H{
			"code": 80000,
			"msg":  "用户id,不存在",
		})
		return
	}
	//info,_ := user.GetUserInfo(userId)
	//if info.Id <= 0 {
	//	c.JSON(200, gin.H{
	//		"code": 80000,
	//		"msg":  "用户id,查询不到",
	//	})
	//	return
	//}
	sendJson.UserId = userId
	//转为Json
	waitSend, _ := json.Marshal(sendJson)
	e := mqUserLevelUser.Publish(viper.GetString("rabbitmqueue.monitorUser"),string(waitSend))
	if e != nil {
		c.JSON(200, gin.H{
			"code": 80000,
			"msg":  e.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 10000,
		"msg": "发送成功",
	})
	return
}
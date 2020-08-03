package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"mlive/library/logger"
	"mlive/library/snowflakeId"
	"mlive/service/admin"
	"mlive/service/java"
	"mlive/service/user/arango"
	"mlive/service/user/dao"
	"mlive/util"
	"strconv"
	"sync"
	"time"
)

var (
	user 		dao.User
	ago         arango.ArangoUser
)

var (
	mqUser           = new(util.RabbitMQ)
)
type MqUserInfo struct {
	UserID 			int64 `json:"userId"`
}

//消费用户消息
func MqUserInfoReceive()  {
	mqUser.Consume(viper.GetString("rabbitmqueue.monitorUser"),mqUserFun)
}
// 更新用户
func mqUserFun(message [] byte) (bool) {
	mqJson := MqUserInfo{}
	json.Unmarshal(message, &mqJson)
	logger.Iprintln(viper.GetString("rabbitmqueue.monitorUser"),"队列接收到消息",mqJson)
	log.Println(viper.GetString("rabbitmqueue.monitorUser"),"队列接收到消息",mqJson)
	//处理你们的业务。。。。
	userId 		 := mqJson.UserID
	no := snowflakeId.GetSnowflakeId()
	user.CreateUserMqlog(userId,0,no)
	var channel   int64 = 0
	var status    int64 = -1
	var msg       string = ""
	javaInfo,err :=java.InfoById(userId)
	if err != nil {
		msg = fmt.Sprintf(" java request  error :(%v) ",  err)
		user.SaveUserMqlog(userId,no, map[string]interface{}{
			"channel":3,
			"status":status,
			"msg":msg,
			"update_time":time.Now(),
		})
		logger.Eprintf(msg)
		return false
	}
	if userId > 0 && javaInfo.UserId >0{
		info,_ := user.GetUserInfo(userId)
		if info.Id > 0 {
			// 更新用户信息
			channel = 2
			newUserId,err :=	user.UpdateSyncUserInfo(javaInfo)
			if err != nil || newUserId <= 0 {
				msg = fmt.Sprintf(" userId : (%d) ,update error : (%v)",userId,err)
				user.SaveUserMqlog(userId,no, map[string]interface{}{
					"channel":channel,
					"status":status,
					"msg":msg,
					"update_time":time.Now(),
				})
				logger.Eprintf(msg)
				return false
			}
		}else{
			channel = 1
			// 注册新用户
			newUserId,err := user.Register(javaInfo)
			if err != nil || newUserId <= 0 {
				msg = fmt.Sprintf(" userId : (%d) ,register error : (%v)",userId,err)
				user.SaveUserMqlog(userId,no, map[string]interface{}{
					"channel":channel,
					"status":status,
					"msg":msg,
					"update_time":time.Now(),
				})
				logger.Eprintf(msg)
				return false
			}

		}
	}else{
		msg = fmt.Sprintf(" java request  error userId:(%d)",userId)
		user.SaveUserMqlog(userId,no, map[string]interface{}{
			"channel":3,
			"status":status,
			"msg":msg,
			"update_time":time.Now(),
		})
		logger.Eprintf(" java request  error userId:(%d)",userId)
		return false
	}

	status = 1
	user.SaveUserMqlog(userId,no, map[string]interface{}{
		"channel":channel,
		"status":status,
		"update_time":time.Now(),
	})

	return true
}


type UserListStruct struct {
	Id  				int64		`json:"userId"`
	Nickname  			string		`json:"nickname"`
	InviteCode  		string		`json:"inviteCode"`
	Mobile				string		`json:"mobile"`
	Level				int64		`json:"level"`
	InviteUpId			int64		`json:"inviteUpId"`
	InviteUpCode		string		`json:"inviteUpCode"`
	InviteUpNickname	string		`json:"inviteUpNickname"`
	MoneyUpId			int64		`json:"moneyUpId"`
	MoneyUpCode		    string		`json:"moneyUpCode"`
	MoneyUpNickname		string		`json:"moneyUpNickname"`
	CreateTime			string		`json:"createTime"`
	UpdateTime			string		`json:"updateTime"`
	Status				int64		`json:"status"`
	Operator			string		`json:"operator"`
}

// 用户列表
func GetUserList(c *gin.Context)  {


	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}
	adminInfo			  ,_:= admin.CheckAdminLogin(adminToken)
	var adminUserId int64 = int64(adminInfo.Id)
	if  adminUserId <= 0 {
		util.Fail(c, 11111, "get admin userId error")
		return
	}

	userId ,_ 		:= strconv.ParseInt(c.DefaultQuery("userId","0"),10,0)
	mobile	  		:= c.DefaultQuery("mobile","")
	nickname  		:= c.DefaultQuery("nickname","")
	level ,_ 		:= strconv.ParseInt(c.DefaultQuery("level","-1"),10,0)
	inviteCode  	:= c.DefaultQuery("inviteCode","")
	inviteUpCode    := c.DefaultQuery("inviteUpCode","")
	inviteUpId,_    :=  strconv.ParseInt(c.DefaultQuery("inviteUpId","0"),10,0)
	moneyUpCode  	:= c.DefaultQuery("moneyUpCode","")
	moneyUpId,_    :=  strconv.ParseInt(c.DefaultQuery("moneyUpId","0"),10,0)
	row ,_    		:= strconv.Atoi(c.DefaultQuery("row","10"))
	page ,_   		:= strconv.Atoi(c.DefaultQuery("page","1"))
	code	  		:= 10000
	msg  	  		:= "success"


	mapBind := map[string]interface{}{
	}
	if userId > 0 {
		mapBind["userId"] = userId
	}

	if mobile != ""{
		mapBind["mobile"] = mobile
	}

	if nickname != ""{
		mapBind["nickname"] = nickname
	}
	if level >= 0 {
		mapBind["level"] = level
	}
	if inviteCode !="" {
		info,_ := user.GetUserCodeInfo(inviteCode)
		if info.Id <= 0 {
			util.Fail(c, 80000, "invite  code errors")
			return
		}
		mapBind["userId"] = info.Id
	}
	if inviteUpId > 0 {
		mapBind["inviteUpId"] = inviteUpId
	}

	if inviteUpCode !="" {
		inUpfo,_ := user.GetUserCodeInfo(inviteUpCode)
		if inUpfo.Id <= 0 {
			util.Fail(c, 80000, "invite up code errors")
			return
		}
		mapBind["inviteUpId"] = inUpfo.Id
	}

	// 分润上级
	// 1、找出所有的下级链条
	// 2、循环判断是否属于上级
	// 3、排序

	if moneyUpCode != "" || moneyUpId > 0{

		if moneyUpCode != "" {
			moneyInfo, _ := user.GetUserCodeInfo(moneyUpCode)
			if moneyInfo.Id <= 0 {
				util.Fail(c, 80000, " money  code errors")
				return
			}
			moneyUpId = moneyInfo.Id
		}
		list,count,_ := ago.GetAdminMoneyList(moneyUpId,mapBind,row, page)
		var data []UserListStruct
		if len(list) > 0 {
			for _,v := range list {

				var (
					inviteUpCode string = ""
					inviteUpNickname string = ""
				)
				upInfo,_ :=user.GetUserInfo(v.InviteUpId)
				if upInfo.Id >0 {
					inviteUpCode 		= upInfo.InviteCode
					inviteUpNickname 	= upInfo.Nickname
				}

				var (
					moneyUpId			int64  = 0
					moneyUpCode		    string = ""
					moneyUpNickname		string = ""

				)
				moneyInfo,_ := ago.GetMoneySuperior(v.Id)
				if moneyInfo.Id > 0 {
					moneyUpId = moneyInfo.Id
					moneyUpCode = moneyInfo.InviteCode
					moneyUpNickname = moneyInfo.Nickname
				}

				data = append(data, UserListStruct{
					Id:               v.Id,
					Nickname:         v.Nickname,
					InviteCode:       v.InviteCode,
					Mobile:           v.Mobile,
					Level:            v.Level,
					InviteUpId:       v.InviteUpId,
					InviteUpCode:     inviteUpCode,
					InviteUpNickname: inviteUpNickname,
					MoneyUpId:        moneyUpId,
					MoneyUpCode:      moneyUpCode,
					MoneyUpNickname:  moneyUpNickname,
					CreateTime:       v.CreateTime,
					UpdateTime:       v.UpdateTime,
					Status:           v.Status,
					Operator:         v.Operator,
				})
			}
		}

		c.JSON(200, gin.H{
			"code": code,
			"msg":  msg,
			"data": map[string]interface{}{
				"items": data,
				"total": count,
			},
		})
	}else {
		list,count,err := ago.GetAdminUserList(mapBind, row, page)
		if err != nil {
			msg = "get user list fail:" + fmt.Sprintf("%v", err)
			util.Fail(c, 80000, msg)
			return
		}

		var data []UserListStruct
		for _,v := range list {

			var (
				inviteUpCode string
				inviteUpNickname string
			)
			upInfo,_ :=user.GetUserInfo(v.InviteUpId)
			if upInfo.Id >0 {
				inviteUpCode 		= upInfo.InviteCode
				inviteUpNickname 	= upInfo.Nickname
			}
			var (
				moneyUpId			int64
				moneyUpCode		    string
		  		moneyUpNickname		string

			)
			moneyInfo,_ := ago.GetMoneySuperior(v.Id)
			if moneyInfo.Id > 0 {
				moneyUpId = moneyInfo.Id
				moneyUpCode = moneyInfo.InviteCode
				moneyUpNickname = moneyInfo.Nickname
			}
			data = append(data,UserListStruct{
				Id:v.Id,
				Nickname:v.Nickname,
				InviteCode:v.InviteCode,
				Mobile:v.Mobile,
				Level:v.Level,
				InviteUpId:v.InviteUpId,
				InviteUpCode:inviteUpCode,
				InviteUpNickname:inviteUpNickname,
				MoneyUpId:moneyUpId,
				MoneyUpCode:moneyUpCode,
				MoneyUpNickname:moneyUpNickname,
				CreateTime:v.CreateTime,
				UpdateTime:v.UpdateTime,
				Status:v.Status,
				Operator:v.Operator,
			})
		}

		c.JSON(200, gin.H{
			"code": code,
			"msg":  msg,
			"data": map[string]interface{}{
				"items": data,
				"total": count,
			},
		})
	}
}



// 用户账号
type UserAccountStruct struct {
	Id  				int64		`json:"userId"`
	Nickname  			string		`json:"nickname"`
	Mobile				string		`json:"mobile"`
	Level				int64		`json:"level"`
	TeamNum   			int64		`json:"teamNum"`
	CloudNum  			int64 		`json:"cloudNum"`
	Remarks  			string 		`json:"remarks"`
	TotalBuy 		    int64 		`json:"totalBuy"`  //累计订购量
	TotalSale 		  	float64 	`json:"totalSale"`  //累计销售
	TotalIncome 	  	float64 	`json:"totalIncome"`  //累计收益
	TotalUavaIncome   	float64 	`json:"totalUavaIncome"`  //在途收益
	WithdrawIncome    	float64 	`json:"withdrawIncome"`  //可提收益
	WithdrawAlready   	float64 	`json:"withdrawAlready"`  //已提现
	OperatorTime		string		`json:"operatorTime"`
	Operator			string		`json:"operator"`
}

func GetUserAccountList(c *gin.Context)  {

	userId ,_ 		:= strconv.ParseInt(c.DefaultQuery("userId","0"),10,0)
	mobile	  		:= c.DefaultQuery("mobile","")
	nickname  		:= c.DefaultQuery("nickname","")
	row ,_    		:= strconv.Atoi(c.DefaultQuery("row","10"))
	page ,_   		:= strconv.Atoi(c.DefaultQuery("page","1"))
	code	  		:= 10000
	msg  	  		:= "success"

	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}



	mapBind := map[string]interface{}{
	}
	if userId > 0 {
		mapBind["userId"] = userId
	}

	if mobile != ""{
		mapBind["mobile"] = mobile
	}

	if nickname != ""{
		mapBind["nickname"] = nickname
	}
	list,count,err := ago.GetAdminUserList(mapBind, row, page)
	if err != nil {
		msg = "get user list fail:" + fmt.Sprintf("%v", err)
		util.Fail(c, 80000, msg)
		return
	}
	var data []UserAccountStruct
	if len(list) > 0 {
		for _,v := range list {
			teamNum,_ := ago.GetTeamNum(v.Id,v.Level)
			cloudNum,_ := cloud.GetCloudStorageSurplusNum(v.Id)
			lastCloudInfo,_ := cloud.GetAdminCloudLastInfo(v.Id)
			var (
				remarks  			string = ""
				operatorTime		string		= ""
				operator			string		= ""
			)
			if lastCloudInfo.Id > 0 {
				remarks = lastCloudInfo.Remarks
				operatorTime = lastCloudInfo.CreateTime.Format("2006-01-02 15:04:05")
				if lastCloudInfo.AdminUserId > 0 {
					adminInfo ,_:= admin.GetAdminUserInfo(adminToken,int(lastCloudInfo.AdminUserId))
					if adminInfo.Id > 0 {
						operator = adminInfo.Username
					}
				}
			}

			var(
				totalBuy 			int64 = 0
				totalSale 			float64 = 0
				totalIncome 		float64 = 0
				totalUavaIncome 	float64 = 0
				withdrawIncome 		float64 = 0
				withdrawAlready 	float64 = 0
			)
			adminAccount,_ := java.GetAdminUserAccountInfoStatistics(v.Id,adminToken)
			if adminAccount.Data.UserId > 0 {
				totalBuy = adminAccount.Data.TotalBuy
				totalSale = adminAccount.Data.TotalSale
				totalIncome = adminAccount.Data.TotalIncome
				totalUavaIncome = adminAccount.Data.TotalUavaIncome
				withdrawIncome = adminAccount.Data.WithdrawIncome
				withdrawAlready = adminAccount.Data.WithdrawAlready
			}

			data = append(data,UserAccountStruct{
				Id:v.Id,
				Nickname:v.Nickname,
				Mobile:v.Mobile,
				Level:v.Level,
				TeamNum:teamNum,
				CloudNum:cloudNum,
				TotalBuy:totalBuy,
				TotalSale:totalSale,
				TotalIncome:totalIncome,
				TotalUavaIncome:totalUavaIncome,
				WithdrawIncome:withdrawIncome,
				WithdrawAlready:withdrawAlready,
				Remarks:remarks,
				OperatorTime:operatorTime,
				Operator:operator,
			})
		}
	}
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": map[string]interface{}{
			"items": data,
			"total": count,
		},
	})
}


// 返回的个人信息
type UserInfo struct {
	UserId    			int64      	`json:"userId" `
	UserName  			string		`json:"userName" `
	Nickname  			string		`json:"nickname" `
	Level     			int64		`json:"level" `
	InviteUpId  		int64		`json:"inviteUpId" `
	TeamNum   			int64		`json:"teamNum" `
	InviteNum 			int64		`json:"inviteNum" `
	MoneyUpId 			int64  		`json:"moneyUpId" `
	Mobile	  		    string		`json:"mobile"`
	CreateTime			string		`json:"createTime"`
	UpdateTime			string		`json:"updateTime"`
	Status				int64		`json:"status"`
	Operator			string		`json:"operator"`

}
/**
* 获取个人信息
 */
func GetUserInfo(c *gin.Context)  {



	userId ,_ := strconv.ParseInt(c.DefaultQuery("userId","0"),10,0)
	code	  := 10000
	msg  	  := "success"
	var userInfo UserInfo
	var mapData map[string]interface{}
	if userId <= 0 {
		util.Fail(c, 11111, "userId parameter fail")
		return
	}
	data, _ := user.GetUserInfo(userId)
	if data.Id > 0 {
		// 团队人数
		userInfo.UserId	   = data.Id
		userInfo.UserName  = data.UserName
		userInfo.Nickname  = data.Nickname
		userInfo.Mobile    = data.Mobile
		userInfo.Level     = data.Level
		userInfo.CreateTime	= data.CreateTime.Format("2006-01-02 15:04:05")
		userInfo.UpdateTime	= data.UpdateTime.Format("2006-01-02 15:04:05")
		userInfo.Status		= data.Status
		userInfo.Operator	= data.Operator
		userInfo.InviteUpId  = data.InviteUpId
		userInfo.TeamNum, _ = ago.GetTeamNum(userId,data.Level)
		// 直属(邀请人数)人数
		userInfo.InviteNum, _ = user.GetUserCount(0, userId)
		agoInfo,_ := ago.GetMoneySuperior(userId)
		if agoInfo.Id > 0 {
			userInfo.MoneyUpId = agoInfo.Id
		}
		mapData = gin.H{
			"code": code,
			"msg":  msg,
			"data": userInfo,
		}
	}else {
		mapData = gin.H{
			"code": code,
			"msg":  msg,
			"data":"",
		}
	}

	c.JSON(200, mapData)
}




// 身份修改
func SaveUserLevel(c *gin.Context)  {

	checkParams   := []interface{}{"userId","level"}
	params,err 	  := util.GetRawData(c.Request,checkParams)
	code	      := 10000
	msg  	      := "success"
	if err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c, 11111, msg)
		return
	}

	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}
	adminInfo			  ,_:= admin.CheckAdminLogin(adminToken)
	var adminUserId int64 = int64(adminInfo.Id)
	if  adminUserId <= 0 {
		util.Fail(c, 11111, "get admin userId error")
		return
	}
	operator := adminInfo.Username
	userParams := util.ChectIntFloat(params["userId"])
	levelParams := util.ChectIntFloat(params["level"])
	if !userParams || !levelParams {
		msg = " parameter fail"
		util.Fail(c, 11111, msg)
		return
	}
	userId := int64(params["userId"].(float64))
	if userId == 1 {
		util.Fail(c, 11111, "user_id == 1 ,不可以操作")
		return
	}
	level := int64(params["level"].(float64))

	info,_ := user.GetUserInfo(userId)
	if info.Id <= 0 {
		util.Fail(c, 80000, " find user error")
		return
	}

	status,msg := user.SaveLevel(userId,0,3,level,adminUserId,operator,"")
	if status  {
		code = 10000

	}else{
		code = 80000
	}
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data":"",
	})
}

// 团队迁移
func SaveUserTeam(c *gin.Context)  {

	checkParams   := []interface{}{"userId","newInviteId"}
	params,err    := util.GetRawData(c.Request,checkParams)
	code	      := 10000
	msg  	      := "success"
	if err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c, 11111, msg)
		return
	}
	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}
	adminInfo			  ,_:= admin.CheckAdminLogin(adminToken)
	var adminUserId int64 = int64(adminInfo.Id)
	if  adminUserId <= 0 {
		util.Fail(c, 11111, "get admin userId error")
		return
	}
	operator := adminInfo.Username

	userParams 		  := util.ChectIntFloat(params["userId"])
	newInviteIdParams := util.ChectIntFloat(params["newInviteId"])
	if !userParams || !newInviteIdParams {
		msg = " parameter fail"
		util.Fail(c, 11111, msg)
		return
	}
	userId 		 := int64(params["userId"].(float64))
	newInviteId  := int64(params["newInviteId"].(float64))

	if userId == 1 {
		msg = " company account cannot be modified"
		util.Fail(c, 80000, msg)
		return
	}
	if userId == newInviteId {
		util.Fail(c,80000,"users can't be the same")
		return
	}

	code ,msg 	 =user.SaveTeam(userId,newInviteId,adminUserId,operator)
	if code == 10000 {
		_,err := user.MqJavaPushUserUpdate(userId)
		if err != nil {
			code = 80000
			msg = " save update user team error"
			logger.Eprintf("save update user id:(%d) team error:(%v)",userId,err)
		}
	}

	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data":"",
	})
}

// 修改用户状态
func SaveStatus(c *gin.Context)  {
	// status 1:正常 2:冻结/黑名单(封禁)
	checkParams   := []interface{}{"userId","status"}
	params,err    := util.GetRawData(c.Request,checkParams)
	code	      := 10000
	msg  	      := "success"

	if err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c, 11111, msg)
		return
	}
	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}
	adminInfo			  ,_:= admin.CheckAdminLogin(adminToken)
	var adminUserId int64 = int64(adminInfo.Id)
	if  adminUserId <= 0 {
		util.Fail(c, 11111, "get admin userId error")
		return
	}
	operator 		  := adminInfo.Username
	userParams 		  := util.ChectIntFloat(params["userId"])
	statusParams := util.ChectIntFloat(params["status"])
	if !userParams || !statusParams {
		util.Fail(c, 11111, "arameter fail")
		return
	}
	userId 		 := int64(params["userId"].(float64))
	if userId == 1 {
		util.Fail(c, 11111, "user_id == 1 ,不可以操作")
		return
	}
	status 		 := int64(params["status"].(float64))
	if status > 2 || status <= 0 {
		util.Fail(c, 80000, "status num fail")
		return
	}

	info,err := user.GetUserInfo(userId)
	if info.Id <=0 || err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c,80000,msg)
		return
	}
	if info.Status == status {
		util.Fail(c,80000,"status equally ")
		return
	}
	mapData := map[string]interface{}{
		"id":userId,
		"status":status,
		"admin_user_id":adminUserId,
		"operator":operator,
	}
	mapStr := map[string]interface{}{
		"status":status,
		"old_status":info.Status,
		"old_admin_user_id":info.AdminUserId,
		"old_update_time":info.UpdateTime.Format("2006-01-02 15:04:05"),
	}
	jsonStr,_ := json.Marshal(mapStr)
	user.CreateUserLog(userId,2,1,adminUserId,string(jsonStr))
	res,err := user.SaveUserInfo(userId,mapData)
	if !res || err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c, 80000, msg)
		return
	}
	// 通过ago
	_,err = ago.Update(userId,mapData)
	if  err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c, 80000, msg)
		return
	}
	// 推送java
	if code == 10000 {
		_,err := user.MqJavaPushUserUpdate(userId)
		if err != nil {
			code = 80000
			msg = " save update user status error"
			logger.Eprintf("save update user id:(%d) status error:(%v)",userId,err)
		}
	}
	util.Fail(c, code, msg)
	return
}

type AdminTeamList struct {
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
	ChildList 	[]*AdminTeamList	`json:"children"`
}


// 团队结构
func GetUserTeamList(c *gin.Context)  {

	userId ,_   	:= strconv.ParseInt(c.DefaultQuery("userId","0"),10,0)
	inviteUpId ,_   := strconv.ParseInt(c.DefaultQuery("inviteUpId","-1"),10,0)
	mobile	  		:= c.DefaultQuery("mobile","")
	nickname  		:= c.DefaultQuery("nickname","")
	if inviteUpId <= 0 {
		inviteUpId = -1
	}

	code := 10000
	msg  := "success"
	var (
		findUid   int64 = 1
	)
	var list []arango.MyUserList
	var teamList []AdminTeamList
	var err  error
	// 顶点

	if mobile != "" {
		info,_ := user.GetUserMobileInfo(mobile)
		if info.Id > 0 {
			userId = info.Id
		}else{
			userId = 999999999999999
		}
	}
	if nickname != "" && userId != 999999999999999 {
		info,_ := user.GetUserNicknameInfo(nickname)

		if userId > 0 && userId<999999999999999 {
			if userId == info.Id {
				userId = info.Id
			}else{
				userId = 999999999999999
			}
		}else{
			if info.Id > 0  {
				userId = info.Id
			}else{
				userId = 999999999999999
			}
		}
	}


	list,err = ago.GetInviteList(userId,inviteUpId)

	if  len(list) >= 1 {
		for _,v := range list {
			//teamNum, _ := ago.GetTeamNum(v.Id,v.Level)
			//// 直属(邀请人数)人数
			//inviteNum, _ := user.GetUserCount(0, v.Id)
			teamList = append(teamList,AdminTeamList{
				UserId:v.Id,
				Nickname:v.Nickname,
				Level:v.Level,
				InviteUpId:v.InviteUpId,
				InviteCode:v.InviteCode,
				TeamNum:0,
				InviteNum:0,
				Status:v.Status,
				ChildList: []*AdminTeamList{},
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
	}

	if  len(teamList) == 1 && 	 inviteUpId == -1 {
		var outList []AdminTeamList
		outId := teamList[0].UserId
		outListAgo,_ := ago.GetInviteOutBoundList(outId)
		if len(outListAgo) > 0 {
			for _,o := range outListAgo {
				//teamNum, _ := ago.GetTeamNum(o.Id,o.Level)
				//// 直属(邀请人数)人数
				//inviteNum, _ := user.GetUserCount(0, o.Id)
				outList = append(outList,AdminTeamList{
					UserId:o.Id,
					Nickname:o.Nickname,
					Level:o.Level,
					InviteUpId:o.InviteUpId,
					InviteCode:o.InviteCode,
					TeamNum:0,
					InviteNum:0,
					Status:o.Status,
					CreateTime:o.CreateTime,
					UpdateTime:o.UpdateTime,
				})
			}
		}
		if len(outList) > 0 {
			wg := &sync.WaitGroup{}
			for i := 0; i < len(outList); i++ {
				wg.Add(1)
				go func(wg *sync.WaitGroup, i int) {
					// 团队人数
					outList[i].TeamNum, _ = ago.GetTeamNum(outList[i].UserId, outList[i].Level)
					//// 直属(邀请人数)人数
					outList[i].InviteNum, _ = user.GetUserCount(0, outList[i].UserId)
					wg.Done()
				}(wg,i)
			}
			wg.Wait()
		}

		findData := make([]*AdminTeamList, 0)  //存储所有初始化struct
	toFindOut:
		if len(outList) > 0 {
			for i,k :=range outList{
				if(k.InviteUpId == findUid){
					findData = append(findData,&k)
					outList = append(outList[:i], outList[i+1:]...)
					findUid = k.UserId
					break
				}
			}
			if len(outList) > 0{
				goto toFindOut;
			}
		}
		findData = append(findData,&teamList[0])
		node := findData[0] //父节点
		makeTree(findData,node)
		var data []*AdminTeamList
		data = append(data,node)
		for i ,v:= range data{
			if v.ChildList == nil {
				data[i].ChildList = []*AdminTeamList{}
			}
		}
		if len(data) == 0 {
			data =  []*AdminTeamList{}
		}


		c.JSON(200, gin.H{
			"code": code,
			"msg":  msg,
			"data":map[string]interface{}{
				"items":data,
				"total":len(data),
			},
		})
	}else{

		if err != nil {
			c.JSON(200, gin.H{
				"code": 80001,
				"msg":  "fail"+fmt.Sprintf("%s",err),
				"data": "",
			})
		} else {




			c.JSON(200, gin.H{
				"code": code,
				"msg":  msg,
				"data":map[string]interface{}{
					"items":teamList,
					"total":len(teamList),
				},
			})
		}
	}

}

// 结构树处理
func makeTree(Allnode []*AdminTeamList,node *AdminTeamList)  {
	childs,_ := haveChild(Allnode,node)
	if childs != nil {
		for _ ,v:= range childs{
			if v.ChildList == nil {
				v.ChildList = []*AdminTeamList{}
			}
		}
		node.ChildList = append(node.ChildList,childs[0:]...) // 添加子节点
		for _,v := range childs{
			_,has := haveChild(Allnode,v) // 查询子节点,判断
			if has  {
				makeTree(Allnode,v) // 递归添加节点
			}
		}
	}
}

// 结构树节点处理
func haveChild(Allnode []*AdminTeamList,node *AdminTeamList)(childs []*AdminTeamList,yes  bool)  {

	for _,v := range Allnode{
		if v.UserId == 1 {
			continue
		}
		if v.InviteUpId == node.UserId {
			childs = append(childs,v)
		}
	}
	if childs != nil {
		yes = true
	}
	return
}

type UserLevelLogData struct {
	 Id	 			int64				`json:"id"`
	 UserId 		int64				`json:"userId"`
	 Nickname 		string				`json:"nickname"`
	 Mobile 		string				`json:"mobile"`
	 OldLevel 		int64				`json:"oldLevel"`
	 NewLevel 		int64				`json:"newLevel"`
	 Type    		int64				`json:"type"`
	 AdminUserId	int64				`json:"adminUserId"`
	 Operator		string				`json:"operator"`
	 CreateTime		string				`json:"createTime"`
}

func UserLevelLog(c *gin.Context)  {


	userId ,_   	:= strconv.ParseInt(c.DefaultQuery("userId","0"),10,0)
	mobile	  		:= c.DefaultQuery("mobile","")
	nickname  		:= c.DefaultQuery("nickname","")
	row ,_    		:= strconv.Atoi(c.DefaultQuery("row","10"))
	page ,_   		:= strconv.Atoi(c.DefaultQuery("page","1"))
	code	  		:= 10000
	msg  	  		:= "success"

	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}


	if mobile != "" {
		info,_ := user.GetUserMobileInfo(mobile)
		if info.Id > 0 {
			userId = info.Id
		}else{
			userId = 999999999999999
		}
	}
	if nickname != "" && userId != 999999999999999 {
		info,_ := user.GetUserNicknameInfo(nickname)

		if userId > 0 && userId<999999999999999 {
			if userId == info.Id {
				userId = info.Id
			}else{
				userId = 999999999999999
			}
		}else{
			if info.Id > 0  {
				userId = info.Id
			}else{
				userId = 999999999999999
			}
		}
	}

	var data []*UserLevelLogData
	list,err := user.GetLevelLogList(userId,page,row,"id desc")
	if err != nil {
		util.Fail(c,80000,fmt.Sprintf("%v",err))
	}
	for _,v := range list{
		var(
			nickname string = ""
			operator string	= ""
			mobile	 string	= ""
		)
		adminInfo ,_:= admin.GetAdminUserInfo(adminToken,int(v.AdminUserId))
		if adminInfo.Id > 0 {
			operator = adminInfo.Username
		}
		info,_ := user.GetUserInfo(v.UserId)
		if info.Id > 0 {
			nickname = info.Nickname
			mobile 	 = info.Mobile
		}
		data = append(data,&UserLevelLogData{
			Id:v.Id,
			UserId:v.UserId,
			Nickname:nickname,
			Mobile:mobile,
			OldLevel:v.OldLevel,
			NewLevel:v.NewLevel,
			Type:v.Type,
			AdminUserId:v.AdminUserId,
			Operator:operator,
			CreateTime:v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	count,_ := user.GetLevelLogCount(userId)

	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data":map[string]interface{}{
			"items":data,
			"total":count,
		},
	})
}

// 设置公司
func SetCompanyAccount(c *gin.Context)  {
	//var (
	//	code = 0
	//	msg = "success"
	//)
	//status,msg := user.SaveLevel(1,3,6,0,"","")
	//if status  {
	//	code = 10000
	//}else{
	//	code = 80000
	//}
	//c.JSON(200, gin.H{
	//	"code": code,
	//	"msg":  msg,
	//	"data":"",
	//})
}
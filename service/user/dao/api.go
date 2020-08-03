package dao

import (
	"encoding/json"
	"mlive/service/java"
	"strconv"
	"time"
	"errors"
	"mlive/library/logger"
	"mlive/service/user/model"
	"mlive/service/user/arango"
	d "mlive/dao"
	"fmt"
)


var(
	ago 		arango.ArangoUser
	tTime 		= time.Now()
)



// 注册邀请用户 userId int64,isCompany int64,Level int64,inviteId int64

func (u *User)Register(result java.JavaUserInfoResult) (int64,error) {

	userId,inviteUpId 	 := result.UserId,result.InviteUserId
	var (
		isCompany 	 int64 = 0
		Level		 int64 = result.PtLevel
	)
	if userId == 1 {
		isCompany = 1
	}
	if userId <= 0 {
		return 0, errors.New("缺少注册用户id")
	}
	if  userId > 1 && inviteUpId <= 0 {
		return 0, errors.New("缺少注册邀请上级")
	}

	userInfo,_ := u.GetUserInfo(userId)
	logger.Iprintln("注册新用户 userId:",userId,"userInfo:",userInfo)
	if userInfo.Id > 0 {
		eMsg := "注册新用户 userId:"+strconv.FormatInt(userId,10)+",已经存在"
		logger.Eprintf(eMsg)
		return 0,errors.New(eMsg)
	}
	if userId == inviteUpId {
		eMsg := "注册新用户 userId:"+strconv.FormatInt(userId,10)+",用户id和邀请上级id一样"
		logger.Eprintf(eMsg)
		return 0,errors.New(eMsg)
	}
	var isTeam int64 = 0
	if inviteUpId > 0 {
		inviteInfo, _ := u.GetUserInfo(inviteUpId)
		logger.Iprintln("注册新用户 userId:",userId,",inviteInfo:",inviteInfo)
		if inviteInfo.Id <= 0 {
			eMsg := "注册新用户 userId:"+strconv.FormatInt(userId,10)+",邀请上级"+strconv.FormatInt(inviteUpId,10)+"不存在"
			logger.Eprintf(eMsg)
			return 0, errors.New(eMsg)
		}
		if inviteInfo.Level >= Level {
			isTeam = 1
		}
	}
	startTime,_  := time.ParseInLocation("2006-01-02 15:04:05",result.AddTime,time.Local)
	updateTime,_ := time.ParseInLocation("2006-01-02 15:04:05",result.UpdateTime,time.Local)
	var data = model.MliveUser{
		Id:userId,
		InviteUpId:inviteUpId,
		InviteTime:startTime,
		UserName:result.UserName,
		HeadImgUrl:result.HeadImgUrl,
		Nickname:result.Nickname,
		Name:result.Name,
		Gender:result.Gender,
		Level:Level,
		InviteCode:result.InviteCode,
		AreaCode:result.AreaCode,
		Mobile:result.Mobile,
		TelPhone:result.TelPhone,
		Email:result.Email,
		IsCompany:isCompany,
		IsTeam:isTeam,
		Status:result.Status,
		Operator:result.Operator,
		CreateTime:startTime,
		UpdateTime:updateTime,
	}
	logger.Iprintln("注册新用户 userId:",userId,"mapData:",data)
	if err := d.Db.DB.Create(&data).Error; err !=nil{
		eMsg := "注册新用户 userId:"+strconv.FormatInt(userId,10)+",创建失败:"+err.Error()
		logger.Eprintf(eMsg)
		return 0,err
	}
	// 建立arangodb 关系
	//nowTime 	:=  tTime.Format("2006-01-02 15:04:05")
	mapData := map[string]interface{}{
		"id":userId,
		"invite_up_id":inviteUpId,
		"invite_time":result.AddTime,
		"user_name":result.UserName,
		"head_img_url":result.HeadImgUrl,
		"nickname":result.Nickname,
		"name":result.Name,
		"gender":result.Gender,
		"level":Level,
		"invite_code":result.InviteCode,
		"area_code":result.AreaCode,
		"mobile":result.Mobile,
		"tel_phone":result.TelPhone,
		"email":result.Email,
		"is_company":isCompany,
		"status":result.Status,
		"operator":result.Operator,
		"create_time":result.AddTime,
		"update_time":result.UpdateTime,
	}
	// 注册ago用户信息
	ago.Create(mapData)
	if inviteUpId > 0 {
		// 建立邀请ago关系链
		ago.CreateOf(userId, inviteUpId)
		if isTeam == 1 {
			// 建立团队关系链接
			ago.CreateTeamOf(userId,inviteUpId)
		}
	}

	// 注册新用户首次->标记为公司
	if userId == 1 && Level < 6 {
		u.SaveLevel(1,0,3,6,0,"","")
	}

	return data.Id,nil
}

/**
* 更新用户信息(java同步过来的信息)
 */
func (u *User)UpdateSyncUserInfo(result java.JavaUserInfoResult) (int64,error) {

	userId     := result.UserId
	userInfo,_ := u.GetUserInfo(userId)
	if userInfo.Id <= 0 {
		return 0, errors.New("用户不存在")
	}
	mapData := map[string]interface{}{
		"id":result.UserId,
		"user_name":result.UserName,
		"head_img_url":result.HeadImgUrl,
		"nickname":result.Nickname,
		"name":result.Name,
		"gender":result.Gender,
		"area_code":result.AreaCode,
		"mobile":result.Mobile,
		"tel_phone":result.TelPhone,
		"email":result.Email,
		"invite_code":result.InviteCode,
	}

	javaUpdateTime,_ := time.ParseInLocation("2006-01-02 15:04:05",result.UpdateTime,time.Local)
	if javaUpdateTime.Unix() > userInfo.UpdateTime.Unix() {
		mapData["update_time"] = javaUpdateTime
		mapData["operator"]   = result.Operator
	}

	logger.Iprintln("更新用户 userId:",userId,"更新mapData:",mapData)
	saveRes,err := u.SaveUserInfo(userId,mapData)
	if !saveRes || err != nil {
		eMsg := "注册新用户 userId:"+strconv.FormatInt(userId,10)+",更新失败:"+err.Error()
		logger.Eprintf(eMsg)
		return 0,err
	}
	if javaUpdateTime.Unix() > userInfo.UpdateTime.Unix() {
		mapData["update_time"] = result.UpdateTime
		mapStr := map[string]interface{}{
			"old_info":userInfo,
			"info":mapData,
		}
		jsonStr,_ := json.Marshal(mapStr)
		u.CreateUserLog(userId,3,1,0, string(jsonStr))


	}
	_,err = ago.Update(userId,mapData)
	if err != nil {
		return 0,err
	}
	return userId,nil
}

/**
* 等级修改
* userId:用户id,
* suitId:1 vip:399,2 店长:2500,3 总监:10000
*channel: 1:用户主动使用(上级赠送),2:系统自动抵扣,3:后台操作,
*level:调整等级,
*adminUserId:管理员id,
*operator:管理员名字
 */
func (user *User)SaveLevel(userId int64,suitId int64,channel int64,level int64,adminUserId int64,operator string,orderNo string)(bool,string)  {
	code	    := true
	msg  	    := "success"

	mapParams := map[string]interface{}{
		"userId":userId,
		"channel":channel,
		"level":level,
		"suitId":suitId,
		"adminUserId":adminUserId,
	}
	info, _ := user.GetUserInfo(userId)

	// 被动升级,如果充值1w, 引力等级判断
	//if  channel == 2  &&  suitId == 3 {
	//	// 缺少java
	//	xfylUserInfo,err :=java.GetXFYLUserInfo(info.AreaCode,info.Mobile)
	//	if err != nil   {
	//		code = false
	//		msg = "get xfyl user error"
	//		mapParams["error_msg"] = "get xfyl user error"
	//		jsonStr1,_ 	   := json.Marshal(mapParams)
	//		user.CreateLevelLog(channel,adminUserId,string(jsonStr1)) // 日志记录
	//		return code,msg
	//	}
	//	if xfylUserInfo.Data.UserId >  0 {
	//		mapParams["ptUserId"] = xfylUserInfo.Data.UserId
	//		mapParams["ptLevel"]  = xfylUserInfo.Data.PtLevel
	//		if xfylUserInfo.Data.PtLevel == 5 {
	//			level = 4 // 支付1w,并且引力总裁级别-->升级合伙人
	//		}
	//		if xfylUserInfo.Data.PtLevel == 6 {
	//			level = 5 // 支付1w,并且引力分公司级别-->升级联创
	//		}
	//		mapParams["level"] = level
	//	}
	//}
	jsonStr,_ 	   := json.Marshal(mapParams)
	levelLogId,err := user.CreateLevelLog(channel,adminUserId,string(jsonStr)) // 日志记录



	if info.Id <= 0 {
		code = false
		msg = "user no exist "
		return code,msg
	}

	if info.Level == level {
		code = false
		msg = "user level equally , not update "
		return code,msg
	}
	dataMap := map[string]interface{}{
		"level":         level,
		"update_time":   tTime,
	}
	if adminUserId > 0 {
		dataMap["admin_user_id"] = adminUserId
		dataMap["operator"] 	 = operator
	}
	res, err := user.SaveUserInfo(userId, dataMap)
	if !res || err != nil {
		code = false
		msg = "user update mysql fail "
		return code,msg
	}
	// 修改ago 信息
	dataMap["update_time"] = tTime.Format("2006-01-02 15:04:05")
	_, err = ago.Update(userId, dataMap)
	if err != nil {
		code = false
		msg = "user update  ago info fail "
		return code,msg
	}

	// 升级处理
	inviteInfo, _ := user.GetUserInfo(info.InviteUpId)
	if level-info.Level > 0 {
		if inviteInfo.Level >= info.Level {
			// 上级处理
			if inviteInfo.Id > 0 && level > inviteInfo.Level {
				user.SaveUserInfo(info.Id, map[string]interface{}{
					"is_team": 0,
				})
				teamRes, err := ago.DeleteTeamOf(info.Id)
				if !teamRes && err != nil {
					code = false
					msg = "user delete  team fail "
					return code,msg
				}

			}
		}
		// 下属处理
		handleBranch, _ := user.GetUserTeamList(0, userId, 0, 2, level)
		if len(handleBranch) > 0 {
			for _, k := range handleBranch {
				// 下属-加入团队
				_, err := ago.CreateTeamOf(k.Id, userId)
				if err != nil {
					code = false
					msg = "user  join  team fail "
					return code,msg
				}
				user.SaveUserInfo(k.Id, map[string]interface{}{
					"is_team": 1,
				})

			}
		}
	} else {
		// 降级处理
		// 上级处理
		if inviteInfo.Level < info.Level {
			if level <= inviteInfo.Level {
				_, err :=ago.CreateTeamOf(info.Id, inviteInfo.Id)
				if err != nil {
					code = false
					msg = "user  join  team fail "
					return code,msg
				}
			}
			user.SaveUserInfo(info.Id, map[string]interface{}{
				"is_team": 1,
			})
		}
		//-下属
		handleBranch, _ := user.GetUserTeamList(0, userId, 1, 1, level)
		if len(handleBranch) > 0 {
			for _, k := range handleBranch {
				// 下属-脱离团队
				_, err := ago.DeleteTeamOf(k.Id)
				if err != nil {
					code = false
					msg = "user delete   team fail "
					return code,msg
				}
				user.SaveUserInfo(k.Id, map[string]interface{}{
					"is_team": 0,
				})

			}
		}
	}

	if code {

		user.SaveLevelLog(levelLogId, map[string]interface{}{
			"user_id":       userId,
			"old_level":     info.Level,
			"new_level":     level,
			"status":        1,
			"admin_user_id": adminUserId,
		})
		// 通知java
		_,err := user.MqJavaPushUserUpdate(userId)
		if err != nil {
			return false,err.Error()
		}
	}

	return code,msg
}

/**
* 邀请人修改
 */
func (user *User)SaveTeam(userId int64,newInviteId int64,adminUserId int64,operator string)(int,string)  {

	code	    := 10000
	msg  	    := "success"

	info, _ := user.GetUserInfo(userId)
	newInviteInfo, _ := user.GetUserInfo(newInviteId)
	if info.Id <= 0 {
		code = 80000
		msg = "user no exist "
		return code,msg
	}

	if newInviteInfo.Id <= 0 {
		code = 80000
		msg = "new invite user no exist "
		return code,msg
	}

	if info.InviteUpId == newInviteInfo.Id {
		code = 80000
		msg = "user invite equally , not update"
		return code,msg
	}
	// 1、判断 新的邀请上级是否在我的邀请下级中
	resAgo, err := ago.CheckIsInviteSubordinate(userId, newInviteId)
	if resAgo {
		code = 80000
		msg = " stay invite subordinate,fail"
		return code,msg
	}
	if err != nil {
		code = 80000
		msg = fmt.Sprintf("%s", err)
		return code,msg
	}
	dataMap := map[string]interface{}{
		"invite_up_id": newInviteId,
		"update_time":  tTime,
		"admin_user_id":adminUserId,
		"operator":operator,
		"is_team":1,
	}
	if info.Level > newInviteInfo.Level {
		dataMap["is_team"] = 0
	}
	// 加入操作日志
	mapStr   := map[string]interface{}{
		"invite_up_id": newInviteId,
		"old_invite_up_id": info.InviteUpId,
		"old_update_time":info.UpdateTime.Format("2006-01-02 15:04:05"),
		"old_admin_user_id":info.AdminUserId,
		"old_is_team":info.IsTeam,
	}
	jsonStr,_ := json.Marshal(mapStr)
	user.CreateUserLog(userId,1,1,adminUserId,string(jsonStr))
	res, err := user.SaveUserInfo(userId, dataMap)
	if !res || err != nil {
		code = 80000
		msg = "user update mysql fail "
		return code,msg
	}
	// 修改ago信息
	dataMap["update_time"] = tTime.Format("2006-01-02 15:04:05")
	_,err = ago.Update(userId, dataMap)
	if err != nil {
		code = 80000
		msg = "user update  ago info fail "
	}
	// 修改邀请关系链
	updateNewinviteId := "mlive_user/" + strconv.FormatInt(newInviteId, 10)
	ago.UpdateOf(userId, map[string]interface{}{
		"_to": updateNewinviteId,
	})
	// 修改团队关系链
	if info.IsTeam == 1 {
		if dataMap["is_team"] == 1 {
			ago.UpdateTeamOf(userId, map[string]interface{}{
				"_to": updateNewinviteId,
			})
		} else {
			ago.DeleteTeamOf(userId)
		}
	} else {
		if dataMap["is_team"] == 1 {
			ago.CreateTeamOf(userId, newInviteId)
		}
	}
	return code,msg
}


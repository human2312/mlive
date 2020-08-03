package dao

import (
	"log"
	"mlive/service/java"
	"strconv"
	"errors"
	"sync"
	d "mlive/dao"
	"mlive/service/user/model"
	"time"
)

/**
* 初始化数据 ago
 */
func (u *User)InitAgoData()(bool,error)  {


	num    := 1
	userVal  := []string{"0"}
	userMergeVal := []string{}
toPage:
	log.Println("userVal:",userVal)
	list,err := u.getInviteList(userVal)
	if err != nil {
		return false,err
	}

	if len(list) > 0 {
		for _, result := range list {
			log.Println("num:",num)
			//log.Println("id:",result.Id)
			num++
			info,err := ago.Info(result.Id)
			if err !=nil {
				return false,err
			}
			startTime 	 :=  result.CreateTime.Format("2006-01-02 15:04:05")
			updateTime   := result.UpdateTime.Format("2006-01-02 15:04:05")
			mapData := map[string]interface{}{
				"id":result.Id,
				"invite_up_id":result.InviteUpId,
				"invite_time":startTime,
				"user_name":result.UserName,
				"head_img_url":result.HeadImgUrl,
				"nickname":result.Nickname,
				"name":result.Name,
				"gender":result.Gender,
				"level":result.Level,
				"mobile":result.Mobile,
				"invite_code":result.InviteCode,
				"tel_phone":result.TelPhone,
				"email":result.Email,
				"is_company":result.IsCompany,
				"status":result.Status,
				"operator":result.Operator,
				"create_time":startTime,
				"update_time":updateTime,
			}
			if info.Id <= 0 {
				// 注册ago用户信息
				ago.Create(mapData)
			}else{
				ago.Update(result.Id,mapData)
			}
			var isTeam int64 = 0
			if result.InviteUpId > 0 {
				inviteInfo, _ := u.GetUserInfo(result.InviteUpId)
				if inviteInfo.Id <= 0 {
					return false, errors.New("邀请人不存在")
				}
				if inviteInfo.Level >= result.Level {
					isTeam = 1
				}
			}
			updateMap := map[string]interface{}{
				"_to": "mlive_user/"+strconv.FormatInt(result.InviteUpId,10),
			}
			if result.InviteUpId > 0 {
				// 建立邀请ago关系链
				infoOf,_  := ago.InfoOf(result.Id)
				if infoOf.Key != "" {
					ago.UpdateOf(result.Id,updateMap)
				}else{
					ago.CreateOf(result.Id, result.InviteUpId)
				}
				if isTeam == 1 {
					// 建立团队关系链接
					infoTeamOf,_ := ago.InfoTeamOf(result.Id)
					if infoTeamOf.Key != "" {
						_,err = ago.UpdateTeamOf(result.Id,updateMap)
					}else{
						_,err =ago.CreateTeamOf(result.Id, result.InviteUpId)
					}
				}else{
					ago.DeleteTeamOf(result.Id)
				}
			}else{
				ago.DeleteOf(result.Id)
				ago.DeleteTeamOf(result.Id)
			}
			userMergeVal = append(userMergeVal,strconv.FormatInt(result.Id,10))
		}
		if len(userMergeVal) > 0 {
			userVal  	 = userMergeVal
			userMergeVal = []string{}
			goto toPage
		}
	}
	return true,nil
}

// 修复邀请上级等级0的用户数据
func (u *User)InitInviteData()(bool,error)  {

	useVal := []string{"0"}
	list,err := u.getInviteList(useVal)
	log.Println("err:",err)
	var status bool = true
	for _,val :=range list  {
		log.Println("val:",*val)
		if val.Id >1 && val.InviteUpId == 0 {
			//
			javaInfo,err :=java.InfoById(val.Id)
			log.Println("请求用户接口:", val.Id,javaInfo)
			if err != nil {
				status = false
				log.Println("请求java 接口失败:", err)
				continue
			}
			if javaInfo.UserId >0 && javaInfo.InviteUserId >0{

				inviteInfo, _ := ago.Info(javaInfo.InviteUserId)
				if inviteInfo.Id <= 0 {
					return false, errors.New("邀请人不存在")
				}

				mapData := map[string]interface{}{
					"id":val.Id,
					"level":javaInfo.PtLevel,
					"invite_up_id":javaInfo.InviteUserId,
				}
				saveRes,err := u.SaveUserInfo(val.Id,mapData)
				if !saveRes || err != nil {
					status = false
					log.Println("保持mysql失败:", err)
					continue
				}

				_,err = ago.Update(val.Id,mapData)
				if err != nil {
					status = false
					log.Println("更新arangodb 失败:", err)
					continue
				}
				updateMap := map[string]interface{}{
					"_to": "mlive_user/"+strconv.FormatInt(javaInfo.InviteUserId,10),
				}
				// 建立邀请ago关系链
				infoOf,_  := ago.InfoOf(val.Id)
				if infoOf.Key != "" {
					ago.UpdateOf(val.Id,updateMap)
				}else{
					ago.CreateOf(val.Id, javaInfo.InviteUserId)
				}
			}
			}else{
				if val.Id > 1 {
					status = false
					log.Println("数据异常:", *val)
					continue
				}
			}


		}
	return status,nil
}

// 初始化同步等级
func (u *User)InitSySnLevel()(bool,error)  {

	var (
		page   int = 1
		row    int = 1000
		orderBy = " id asc "
		status = true
		arr  = []int64{}
	)
	wg := &sync.WaitGroup{}
	gTO:
	list,_ := u.GetUserList(0,page,row,orderBy)
		for i := 0; i < len(list); i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, i int) {
				javaInfo, err := java.InfoById(list[i].Id)
				if err != nil {
					log.Println("java request:",list[i].Id,err)
				}
				if javaInfo.UserId > 0 && javaInfo.InviteUserId > 0 {
					if list[i].Level != javaInfo.PtLevel {
						status = false
						arr = append(arr,list[i].Id)
						log.Println("等级对应不上 userId:",list[i].Id, "go-手机:",list[i].Mobile, "go-等级:",list[i].Level, "java-等级:",javaInfo.PtLevel)
					}
				}
				wg.Done()
			}(wg, i)

		}
	wg.Wait()
	if len(list) > 0 {
		page++
		goto gTO;
	}
	log.Println("修改修复的数据arr:",arr)
	return status,nil
}

type UserErrData struct{
	Id 	int `gorm:"column:id"`
	Mobile 	string `gorm:"column:mobile"`
}

// 同步状态
func (u *User)InitSyncStatus()(bool,error)  {

	var data = []*UserErrData{}
	var table string = "mlive_user_mobile_err"
	dbc :=  d.Db.DB
	if err := dbc.Table(table).Where("status = ?",0).Find(&data).Error ; err != nil{
		return  false,errors.New("获取用户失败")
	}
	wg := &sync.WaitGroup{}
	if len(data) > 0 {
		for _,val := range data {

			wg.Add(1)
			go func(wg *sync.WaitGroup,pid int,mobile string) {
					var userData = model.MliveUser{}
					if err := dbc.Table("mlive_user").Where("mobile=?",mobile).First(&userData).Error; err != nil {
						log.Println("用户不存在 手机号码:",mobile)
					}
					if userData.Id <=0 {
						log.Println("用户不存在 userData = 0 手机号码:",mobile)
					}else {
						nowTime 	:=  time.Now()
						log.Println("处理用户:", userData.Id)
						mapData := map[string]interface{}{
							"id":     userData.Id,
							"status": 2,
							"update_time":nowTime,
						}
						if err := dbc.Table("mlive_user").Model(&model.MliveUser{Id:userData.Id}).Updates(mapData).Error; err != nil {
							log.Println("mysql 用户 修改状态失败", userData.Id)
						}else{
							mapData["update_time"] = nowTime.Format("2006-01-02 15:04:05")
							_, err := ago.Update(userData.Id, mapData)
							if err != nil {
								log.Println("ago 用户 修改状态失败:", userData.Id)
							}else{
								mapStatus := map[string]interface{}{
									"status": 1,
								}
								if err := dbc.Table(table).Where("id = ?", pid).Updates(mapStatus).Error; err != nil {
									log.Println("修改user err 失败:", userData.Id)
								}else{
									// mq 通知
									_, err = u.MqJavaPushUserUpdate(userData.Id)
									if err != nil {
										msg := " save update user status error"
										log.Println(msg, userData.Id)
									}
								}
							}
						}
						wg.Done()
					}
				}(wg,val.Id,val.Mobile)

			}
		}
		wg.Wait()

	return true,nil
}

func (u *User)InitSyncMq()(bool,error)  {

	var table string = "mlive_user"
	var data = []*model.MliveUser{}
	dbc :=  d.Db.DB
	if err := dbc.Table(table).Where("status = ?",2).Find(&data).Error ; err != nil{
		return  false,errors.New("获取用户失败")
	}
	if len(data) > 0 {
		for _,val := range data {
				log.Println("data userId:",val.Id)
				if val.Id <=0 {
					log.Println("用户不存在 userData = 0 手机号码:",val.Mobile,val.Id)
				}else{
					// mq 通知
					_, err := u.MqJavaPushUserUpdate(val.Id)
					if err != nil {
						msg := " save update user status error"
						log.Println(msg, val.Id)
					}
				}

		}
	}
	return true,nil
}
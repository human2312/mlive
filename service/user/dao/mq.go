package dao

import (
	"encoding/json"
	"github.com/spf13/viper"
	"mlive/service/user/model"
	"mlive/util"
	"time"
	d "mlive/dao"
)



// 推送-修改的用户状态信息


type PushJavaUserRes struct {
	UserID 			int64 `json:"userId"`
	Level 			int64 `json:"level"`
	InviteUpUserId 	int64 `json:"inviteUpUserId"`
	Status 			int64 `json:"status"`
}
// 给java推送用户
func (u *User)MqJavaPushUserUpdate(userId int64)(bool,error)  {
	info,err := u.GetUserInfo(userId)
	var pushJavaUser           = new(util.RabbitMQ)
	if info.Id > 0 {
		sendJson := PushJavaUserRes{}
		sendJson.UserID = info.Id
		sendJson.Level = info.Level
		sendJson.InviteUpUserId = info.InviteUpId
		sendJson.Status			= info.Status
		//转为Json
		waitSend, _ := json.Marshal(sendJson)
		err := pushJavaUser.Publish(viper.GetString("rabbitmqueue.pushUser"),string(waitSend))
		if err != nil {
			return  false,err
		}
		return true,nil
	}else {
		return false, err
	}
}


// 保持用户mq记录
func (user *User)CreateUserMqlog(userId int64,channel int64,no int64)(int64,error){
	var data = model.MliveUserMqLog{
		UserId:userId,
		Channel:channel,
		No:no,
		Status:0,
		CreateTime:time.Now(),
		UpdateTime:time.Now(),
	}
	if err := d.Db.DB.Create(&data).Error; err !=nil{
		return 0,err
	}
	return data.Id,nil
}

// 修改mq记录
func (user *User)SaveUserMqlog(userId int64,no int64,mapData map[string]interface{})(bool,error)  {
	if err := d.Db.DB.Model(&model.MliveUserMqLog{}).Where("user_id = ?",userId).Where("no = ?",no).Updates(mapData).Error; err !=nil {
		return false,err
	}else{
		return true,nil
	}
}
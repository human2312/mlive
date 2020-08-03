package model

import (
	"time"
)
func (MliveUser)TableName() string {
	return "mlive_user"
}

type MliveUser struct {
	Id 			 int64			`json:"userId" 			gorm:"column:id"`
	InviteUpId	 int64			`json:"inviteUpId"		gorm:"column:invite_up_id"`
	InviteTime   time.Time		`json:"-"				gorm:"column:invite_time"`
	UserName	 string			`json:"userName" 		gorm:"column:user_name"`
	HeadImgUrl	 string			`json:"headImgUrl" 		gorm:"column:head_img_url"`
	Nickname	 string			`json:"nickname" 		gorm:"column:nickname"`
	Name    	 string			`json:"name" 			gorm:"column:name"`
	Gender		 int64			`json:"gender"			gorm:"column:gender"`
	Level		 int64			`json:"level" 			gorm:"column:level"`
	InviteCode	 string			`json:"inviteCode" 		gorm:"column:invite_code"`
	AreaCode	 string			`json:"areaCode"		gorm:"column:area_code"`
	Mobile	 	 string			`json:"mobile"			gorm:"column:mobile"`
	TelPhone	 string			`json:"telPhone"		gorm:"column:tel_phone"`
	Email		 string			`json:"email"			gorm:"column:email"`
	IsCompany  	 int64			`json:"isCompany"		gorm:"column:is_company"`
	IsTeam  	 int64			`json:"isTeam"			gorm:"column:is_team"`
	Status		 int64			`json:"status"			gorm:"column:status"`
	Operator	 string			`json:"operator"		gorm:"column:operator"`
	AdminUserId  int64			`json:"adminUserId"		gorm:"column:admin_user_id"`
	CreateTime   time.Time		`json:"createTime"		gorm:"column:create_time"`
	UpdateTime   time.Time		`json:"updateTime"		gorm:"column:update_time"`
}



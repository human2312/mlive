package model

import (
	"time"
)
func (MliveUserCodeLog)TableName() string {
	return "mlive_user_code_log"
}

type MliveUserCodeLog struct {
	Id 			 int64			`json:"id" 				gorm:"column:id"`
	UserId 		 int64			`json:"userId" 			gorm:"column:user_id"`
	No	 		 int64			`json:"no" 				gorm:"column:no"`
	Type 		 int64			`json:"type" 			gorm:"column:type"`
	CodeType 	 int64			`json:"codeType" 			gorm:"column:code_type"`
	UseUserId 	 int64			`json:"useUserId" 		gorm:"column:use_user_id"`
	Number 	 	 int64			`json:"number" 			gorm:"column:number"`
	Channel 	 int64			`json:"channel" 		gorm:"column:channel"`
	Status  	 int64			`json:"status"		gorm:"column:status"`
	AdminUserId  int64			`json:"adminUserId"		gorm:"column:admin_user_id"`
	CreateTime   time.Time		`json:"-"				gorm:"column:create_time"`
	UpdateTime   time.Time		`json:"-"				gorm:"column:update_time"`
}




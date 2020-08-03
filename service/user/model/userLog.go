package model

import (
	"time"
)

func (MliveUserLog)TableName() string {
	return "mlive_user_log"
}

type MliveUserLog struct {
	Id 			 int64			`json:"userId" 			gorm:"column:id"`
	Type		 int64			`json:"type"			gorm:"column:type"`
	UserId		 int64			`json:"userId"			gorm:"column:user_id"`
	Status	 	 int64			`json:"status"			gorm:"column:status"`
	AdminUserId  int64			`json:"adminUserId"		gorm:"column:admin_user_id"`
	Params  	 string			`json:"params"			gorm:"column:params"`
	CreateTime   time.Time		`json:"-"				gorm:"column:create_time"`
	UpdateTime   time.Time		`json:"-"				gorm:"column:update_time"`
}


package model

import (
	"time"
)


func (MliveUserLevelLog)TableName() string {
	return "mlive_user_level_log"
}

type MliveUserLevelLog struct {
	Id 			 int64			`json:"userId" 			gorm:"column:id"`
	Type		 int64			`json:"type"			gorm:"column:type"`
	UserId		 int64			`json:"userId"			gorm:"column:user_id"`
	OldLevel	 int64			`json:"oldLevel"		gorm:"column:old_level"`
	NewLevel	 int64			`json:"newLevel"		gorm:"column:new_level"`
	Status	 	 int64			`json:"status"			gorm:"column:status"`
	AdminUserId  int64			`json:"adminUserId"		gorm:"column:admin_user_id"`
	Params  	 string			`json:"params"			gorm:"column:params"`
	CreateTime   time.Time		`json:"-"				gorm:"column:create_time"`
	UpdateTime   time.Time		`json:"-"				gorm:"column:update_time"`
}



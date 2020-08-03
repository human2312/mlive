package model

import (
	"time"
)

func (MliveUserMqLog)TableName() string {
	return "mlive_user_mq_log"
}

type MliveUserMqLog struct {
	Id 			 int64			`json:"userId" 			gorm:"column:id"`
	UserId		 int64			`json:"userId"			gorm:"column:user_id"`
	Channel		 int64			`json:"channel"			gorm:"column:channel"`
	No	 		 int64			`json:"no" 				gorm:"column:no"`
	Status	 	 int64			`json:"status"			gorm:"column:status"`
	Msg			 string			`json:"msg"				gorm:"column:msg"`
	CreateTime   time.Time		`json:"-"				gorm:"column:create_time"`
	UpdateTime   time.Time		`json:"-"				gorm:"column:update_time"`
}


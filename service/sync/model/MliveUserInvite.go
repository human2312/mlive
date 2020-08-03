package model

import (
	"time"
)

type MliveUserInvite struct {
	ID 			 int64			`json:"userId" 		gorm:"column:id"`
	Pid 		 int64			`json:"pId" 		gorm:"column:pid"`
	CreateTime   time.Time		`json:"-"			gorm:"column:create_time"`
	UpdateTime   time.Time		`json:"-"			gorm:"column:update_time"`
}
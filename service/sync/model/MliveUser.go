package model

import (
	"time"
)
type MliveUser struct {
	ID 			 int64			`json:"userId" 		gorm:"column:id"`
	XfylId		 int64			`json:"xfylId" 		gorm:"column:xfyl_id"`
	UserName	 string			`json:"userName" 	gorm:"column:user_name"`
	Password	 string			`json:"password" 	gorm:"column:password"`
	HeadImgUrl	 string			`json:"headImgUrl" 	gorm:"column:head_img_url"`
	Nickname	 string			`json:"nickname" 	gorm:"column:nickname"`
	Name		 string			`json:"name" 		gorm:"column:name"`
	Level		 int64			`json:"level" 		gorm:"column:level"`
	Gender		 int64			`json:"gender" 		gorm:"column:gender"`
	Mobile		 string			`json:"mobile" 		gorm:"column:mobile"`
	TelPhone	 string			`json:"telPhone" 	gorm:"column:tel_phone"`
	Status	 	 int64			`json:"status" 		gorm:"column:status"`
	Deleted	 	 int64			`json:"deleted" 	gorm:"column:deleted"`
	Shield	 	 int64			`json:"shield" 		gorm:"column:shield"`
	CreateTime   time.Time		`json:"-"			gorm:"column:create_time"`
	UpdateTime   time.Time		`json:"-"			gorm:"column:update_time"`
}




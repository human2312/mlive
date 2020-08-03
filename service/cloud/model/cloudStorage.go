package model

import (
	"time"
)
func (MliveCloudStorage)TableName() string {
	return "mlive_cloud_storage"
}

type MliveCloudStorage struct {
	Id 			 int64			`json:"id" 				gorm:"column:id"`
	UserId 		 int64			`json:"userId" 			gorm:"column:user_id"`
	No	 		 int64			`json:"no" 				gorm:"column:no"`
	OrderNo	 	 string			`json:"orderNo" 		gorm:"column:order_no"`
	Type 		 int64			`json:"type" 			gorm:"column:type"`
	UseUserId 	 int64			`json:"useUserId" 		gorm:"column:use_user_id"`
	CloudType 	 int64			`json:"cloudType" 		gorm:"column:cloud_type"`
	Number 	 	 int64			`json:"number" 			gorm:"column:number"`
	Channel 	 int64			`json:"channel" 		gorm:"column:channel"`
	Status  	 int64			`json:"status"			gorm:"column:status"`
	Remarks		 string			`json:"remarks"			gorm:"column:remarks"`
	AdminUserId  int64			`json:"adminUserId"		gorm:"column:admin_user_id"`
	CreateTime   time.Time		`json:"-"				gorm:"column:create_time"`
	UpdateTime   time.Time		`json:"-"				gorm:"column:update_time"`
}




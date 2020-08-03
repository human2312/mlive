package model

import (
	"time"
)
func (MliveUserCodeSummary)TableName() string {
	return "mlive_user_code_summary"
}

type MliveUserCodeSummary struct {
	Id 			 	 int64			`json:"id" 				gorm:"column:id"`
	UserId 		 	 int64			`json:"userId" 			gorm:"column:user_id"`
	Goods 		 	 int64			`json:"goods" 			gorm:"column:goods"`
	Vip 		 	 int64			`json:"vip" 			gorm:"column:vip"`
	ShopOwner 		 int64			`json:"shopOwner" 		gorm:"column:shop_owner"`
	ChiefInspector 	 int64			`json:"chiefInspector" 	gorm:"column:chief_inspector"`
	Partner 		 int64			`json:"partner" 		gorm:"column:partner"`
	Remarks			 string			`json:"remarks" 		gorm:"column:remarks"`
	AdminUserId 	 int64			`json:"adminUserId"		gorm:"column:admin_user_id"`
	CreateTime   	 time.Time		`json:"-"				gorm:"column:create_time"`
	UpdateTime   	 time.Time		`json:"-"				gorm:"column:update_time"`
}




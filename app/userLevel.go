package app

import (
	"encoding/json"
	"github.com/spf13/viper"
	"mlive/util"
	//"mlive/service/user/arango"
	"mlive/service/user/dao"
)

var (
	userL 		dao.User
	//agoL        arango.ArangoUser
)
var (
	mqUserLevel           = new(util.RabbitMQ)
)
type MqUserLevelInfo struct {
	UserId 			int64 	`json:"userId"`
	Type 			int64 	`json:"type"`
	No 				string 	`json:"no"`
	NoType 			int 	`json:"noType"`
}

//消费用户消息
func MqUserLevelReceive()  {
	mqUser.Consume(viper.GetString("rabbitmqueue.monitorUserLevel"),mqUserLevelFun)
}

// 更新用户等级
func mqUserLevelFun(message [] byte) (bool) {
	mqJson := MqUserLevelInfo{}
	json.Unmarshal(message, &mqJson)


	//处理你们的业务。。。。
	//userId 		  := mqJson.UserId
	//ty 		 	  := mqJson.Type
	//no 			  := mqJson.No
	//noType 	  	  := mqJson.NoType
	//var level	 int64 = -1
	//info,err := userL.GetUserInfo(userId)
	//if err != nil {
	//	logger.Eprintf("   error userId:(%d)",userId)
	//	return false
	//}
	//if info.Id <= 0 {
	//	logger.Eprintf("   find  userId:(%d) error",userId)
	//	return false
	//}
	//var code bool = false
	//
	//if ty == 1 {
	//	// 1 兑换条码
	//	if noType == 1 || noType == 2 || noType == 3 || noType == 4 {
	//		if noType == 1  {
	//			level = 1
	//		}
	//		if noType == 2  {
	//			level = 2
	//		}
	//		if noType == 3  {
	//			level = 2
	//		}
	//		if noType == 4  {
	//			level = 2
	//		}
	//		if level >= 1  && level  > info.Level  {
	//			code,_ = userL.SaveLevel(userId,1,level,0,"")
	//
	//		}
	//
	//	}
	//	}else if ty == 2 {
	//	// 2 订单号
	//	if noType == -1 { //购买商品99
	//		data,err := java.GetOrderInfo(no)
	//		if err != nil {
	//			logger.Eprintf("   java  GetDepositInfo error :(%v) ",err)
	//			return false
	//		}
	//		if data.Data.OrderInfo.Id > 0 && data.Data.OrderInfo.Status != "W" &&  info.Level == 0 {
	//			code,_ = userL.SaveLevel(userId,1,1,0,"")
	//			return code
	//		}
	//
	//	// recommend-2:赚播店长999,recommend-3:赚播总监9999,recommend-4:赚播合伙人30000,recommend-5:赚播联创89100
	//	}else if noType == 2 || noType == 3 || noType == 4 || noType == 5 {
	//		//port\invest
	//		// 充值订单号
	//		data,err := java.GetDepositInfo(no)
	//		log.Println("data:",data)
	//		if err != nil {
	//			logger.Eprintf("   java  GetDepositInfo error :(%v) ",err)
	//			return false
	//		}
	//		if userId != data.Data.UserId {
	//			logger.Eprintf("   java  GetDepositInfo  userid error :(%d ) : %v ",userId,no)
	//			return false
	//		}
	//
	//		if data.Data.Id > 0 && data.Data.OrderStatus == "S"  {
	//			if noType == 2 && info.Level == 1 {
	//				level = 2
	//			}
	//			if noType == 3 && info.Level == 2 {
	//				level = 3
	//			}
	//			if noType == 4 && info.Level == 3 {
	//				level = 4
	//			}
	//			if noType == 5 && info.Level == 4 {
	//				level = 5
	//			}
	//		}
	//		if level >= 2  && level  > info.Level  {
	//			code,_ = userL.SaveLevel(userId,1,level,0,"")
	//
	//		}
	//	}
	//}
	//
	//
	//return  code
	return true
}

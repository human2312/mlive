package dao

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	d "mlive/dao"
	"mlive/library/snowflakeId"
	"mlive/service/cloud/model"
	"mlive/service/java"
	userDaoc "mlive/service/user/dao"
	"mlive/util"
	"time"
	"fmt"
)
var(
	userc = userDaoc.User{}
)

// 云仓创建
func (cloud *CloudStorage)CreateCloudLog(userId int64,types int64,orderNo string,useUserId int64,cloudType int64,number int64,channel int64,status int64,adminUserId int64,remarks string)(int64,error)  {

	no := snowflakeId.GetSnowflakeId()
	var data = model.MliveCloudStorage{
		UserId:userId,
		No:no,
		OrderNo:orderNo,
		Type:types,
		UseUserId:useUserId,
		CloudType:cloudType,
		Number:number,
		Channel:channel,
		Status:status,
		AdminUserId:adminUserId,
		Remarks:remarks,
		CreateTime:time.Now(),
		UpdateTime:time.Now(),
	}
	if err := d.Db.DB.Create(&data).Error; err !=nil{
		return 0,err
	}
	if  data.Id > 0 {
		return no,nil
	}else{
		return 0,errors.New("create cloud log error")
	}
}

// 修改云仓订单信息
func (cloud *CloudStorage)SaveCloudLogByOrder(userId int64,No int64,mapData map[string]interface{})(bool,error) {
	//mapData["quantity"] = gorm.Expr("quantity-?",1)
	db := d.Db.DB.Model(&model.MliveCloudStorage{})
	if userId > 0 {
		db = db.Where("user_id=?", userId)
	}
	if No > 0 {
		db = db.Where("no =?", No)
	}
	if  No <= 0 {
		return false,errors.New("[no] fail ")
	}
	if err := db.UpdateColumn(mapData).Error; err != nil {
		return false, err
	} else {
		return true, nil
	}
}


// 处理云仓-赠品商品状态
// @cloudNo 云仓编号
func (cloud *CloudStorage)SaveGiveGoodsStatus(cloudNo int64)(bool,error)  {

	info,err := cloud.GetCloudStoragOrderInfo(0,cloudNo)

	if err != nil || info.Id <= 0   {
		return false,err
	}
	if info.Status != 0 {
		return false,errors.New("status fail ")
	}
	res,err := cloud.SaveCloudLogByOrder(0,cloudNo, map[string]interface{}{
		"status":1,
		"update_time":time.Now(),
	})
	return res,err
}

//保存-云仓
// @userId		原消耗库存的人
// @types		1 新增，2 减少
// @cloudType	类型 0:@types 3个(赠送商品)，1:@types21个(vip),2:@types 147(店长),3:@types 588总监,4:@types 1764(合伙人),5:types 5292(联创)
// @useUserId	原userId的下级
// @orderNo		订单号
//channel: 1:用户主动使用(上级赠送) 线下,2:系统自动抵扣 线上,3:后台操作,
func (cloud  *CloudStorage)SaveReduceCloudStoragNums(userId int64,types int64,cloudType int64,useUserId int64,orderNo string,channel int64,adminUserId int64)(int64,error)  {

	// 订单号存在 不扣云仓
	orderInfo, _ := cloud.GetCloudStoragOrderNoInfo(orderNo,types)
	if orderInfo.Id > 0 {
		return 0, errors.New("order no already exist")
	}
	if channel != 2 {
		return 0, errors.New(" cloud channel error ")
	}

	var num int64 = 0
	if cloudType == 0  {
		num = 3
	}else if cloudType == 1 {
		num = 21
	}else if cloudType == 2 { //店长
		num = 29
	}else if cloudType == 3 {
		num = 117
	}else if cloudType == 4 { //1764
		num = 294
	}else if cloudType == 5 { //5292
		num = 17647
	}else{
		return 0,errors.New("cloudType   error ")
	}
	if types == 2 {
		number,_ := cloud.GetCloudStorageSurplusNum(userId)
		if number <=0 {
			return 0,errors.New("number  shortage < 0 ")
		}
		if number < num {
			return 0,errors.New("number < num shortage ")
		}
		num = -num
	}
	if num == 0 {
		return 0,errors.New("num For 0 ")
	}
	if types == 1 && cloudType<=0{
		return 0,errors.New("add cloud error type ")
	}

	// 判断是否复购
	//info,_ := userc.GetUserInfo(userId)
	//if info.Id > 0 {
	//	//   是复购
	//	if cloudType <= info.Level && types == 1 && cloudType >=1 {
	//
	//		if  info.Level == 1 { //vip
	//				num = 22
	//		}else if  info.Level == 2 {//店长
	//				if cloudType == 1 {
	//					num = 29
	//				}else {
	//					num = 147
	//				}
	//		}else if info.Level == 3  { //总监
	//			if cloudType == 1 {
	//				num = 33
	//			}else if cloudType == 2 {
	//				num = 167
	//			} else{
	//				num = 667
	//			}
	//		}else if info.Level == 4  { //合伙人
	//			if cloudType == 1 {
	//				num = 38
	//			}else if cloudType == 2 {
	//				num = 192
	//			}else if cloudType == 3 {
	//				num = 769
	//			} else{
	//				num = 2308
	//			}
	//		}else if info.Level == 5  { //联创
	//			if cloudType == 1 {
	//				num = 42
	//			}else if cloudType == 2 {
	//				num = 208
	//			}else if cloudType == 3 {
	//				num = 833
	//			}else if cloudType == 4 {
	//				num = 2500
	//			}  else{
	//				num = 7500
	//			}
	//		}else{
	//			return false,errors.New("info types error  ")
	//		}
	//	}
	//}else{
	//	return false,errors.New("user error ")
	//}

	no,err := cloud.CreateCloudLog(userId,types,orderNo,useUserId,cloudType,num,channel,1,adminUserId,"")
	if no <= 0 {
		return 0,err
	}
	// 新增云仓通知
	if types == 1  && channel == 2 {
		cloud.MqJavaPushUserUpdate(orderNo,userId,useUserId,num,no,1)
	}

	return num,nil
}

// 购买商品-减少库存
func (cloud  *CloudStorage)SaveReduceCloudStoragGoodsNums(userId int64,useUserId int64,num int64,orderNo string)(int64,error)  {

	var cloudType int64  = 0
	var channel int64 = 2 // 系统自动抵扣:用户主动购买
	var types int64 = 2

	if num == 0 {
		return 0,errors.New("num For 0 ")
	}
	number,_ := cloud.GetCloudStorageSurplusNum(userId)
	if number <=0 {
		return 0,errors.New("number  shortage < 0 ")
	}
	if number < num {
		return 0,errors.New("number < num shortage ")
	}
	no,err := cloud.CreateCloudLog(userId,types,orderNo,useUserId,cloudType,-num,channel,1,0,"")
	if no <= 0 {
		return 0,err
	}
	cloud.MqJavaPushUserUpdate(orderNo,useUserId,userId,num,no,1)
	return num,nil
}



// 推送-扣云仓推送
var (
	pushJavaCloud           = new(util.RabbitMQ)
)

type PushJavaCloudRes struct {
	OrderNo 		string		 `json:"orderNo"`//订单号
	UserId 			int64 		 `json:"userId"` //下单人userId
	ProfitUserId 	int64 		 `json:"profitUserId"` //分润人ID(被扣云仓数量的userId)
	BuyNum 			int64 		 `json:"buyNum"` //购买数量(扣减云仓数量)
	CouponSn 		int64 		 `json:"couponSn"` //云仓编号
	RewardType		int64 		 `json:"rewardType"` // 1:商品(线上) 2:服务商(线下)
}
// 给java推送用户
func (cloud *CloudStorage)MqJavaPushUserUpdate(OrderNo string,UserId int64,ProfitUserId int64,BuyNum int64,couponSn int64,rewardType int64)(bool,error)  {

		sendJson := PushJavaCloudRes{}
		sendJson.OrderNo 		= OrderNo
		sendJson.UserId 		= UserId
		sendJson.ProfitUserId 	= ProfitUserId
		sendJson.BuyNum			= BuyNum
		sendJson.CouponSn		= couponSn
		sendJson.RewardType		= rewardType

		//转为Json
		waitSend, _ := json.Marshal(sendJson)
		pushJavaCloud.Publish(viper.GetString("rabbitmqueue.cloud"),string(waitSend))
		return true,nil

}



type TimingCloudInfo struct {

}
//定时检查-云仓订单是否
func (this TimingCloudInfo)Run() {

	var  cloud CloudStorage
	list,err := cloud.GetCloudStoragTimingList()
	if err != nil {
		fmt.Println(" timing cloud list err:",err.Error())
	}
	if len(list)  > 0 {
		for _,v := range list {
			orderInfo,err := java.GetOrderInfo(v.OrderNo)
			if err != nil {
				fmt.Println("timing cloud java request err GetOrderInfo:",err.Error())
				continue
			}
			fmt.Println("timing cloud order no:",v.OrderNo)
			if orderInfo.Data.OrderInfo.Id > 0 && v.No > 0 {
				//订单状态 未支付/待支付:W 已完成:S 已取消:C 待发货:WS 待收货:WD
				fmt.Println("order status:", orderInfo.Data.OrderInfo.Status)
				if orderInfo.Data.OrderInfo.Status == "C" {
					cloud.SaveCloudLogByOrder(0,v.No, map[string]interface{}{
						"status":-1,
						"update_time":time.Now(),
					})
				}
				if orderInfo.Data.OrderInfo.Status == "WS" || orderInfo.Data.OrderInfo.Status == "WD" || orderInfo.Data.OrderInfo.Status == "S" {
					cloud.SaveCloudLogByOrder(0,v.No, map[string]interface{}{
						"status":1,
						"update_time":time.Now(),
					})
				}
			}
		}
	}
}





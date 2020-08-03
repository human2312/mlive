package java

import (
	"encoding/json"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

// 充值信息
type DepositResult struct {
	Msg					 string     `json:"msg"`
	Code				 int64		`json:"code"`
	Data 				 struct{
		Id 		 	 		 int64 		`json:"id"`  //id
		PayNo				 string 		 `json:"payNo"` //支付流水
		UserId 		 	 	 int64 		`json:"userId"`  //id
		DepositNo 		 	 string 	`json:"depositNo"`  //充值订单号
		Price				 float64    `json:"price"` //订单金额
		OrderStatus		     string    	`json:"orderStatus"` //状态 待充值:W;充值成功:S;充值确认:BW;充值失败:F
		BankName			 string 		 `json:"bankName"`//付款方式
		OpenId				 string 		 `json:"openId"` //付款账号,银行用户标识
		PayDate				 string 	`json:"payDate"` //支付时间-支付日期 yyyy-MM-dd
		PayTime				 string 	`json:"payTime"` //支付时间-支付时间 HH:mm:ss
	}		`json:"data"`
}
// 获取充值信息
func GetDepositInfo(depositNo string)(DepositResult,error)  {
	info := DepositResult{}
	params := map[string]interface{}{
		"depositNo":depositNo,
	}
	res,err  := httpApiPost(viper.GetString("java.java_user_deposit_url"),params,"v1","")
	json.Unmarshal(res, &info)
	return info,err
}

type OrderDepositResult struct {
	Msg					 string     `json:"msg"`
	Code				 int64		`json:"code"`
	Data 				 struct{
		PrePayNo 		 	 string 		`json:"prePayNo"`  //预支付订单号
		OrderNo 		 	 string 		`json:"orderNo"`  //订单号
		MercId 		 	 	 string 		`json:"mercId"`  //商户号
	}
}

// java_order_deposit_submit : 订单提交
// @no 云仓编号
// @ 订单类型  1:VIP499 ，2:赚播店长2500 ，3:赚播总监10000 ，4:赚播合伙人30000 ，5:赚播联创90000
// @ 下单人的下级:userId
func SubmitOrderDeposit(no int64,orderType int64,userId int64,buyNum int64,userToken string)(OrderDepositResult,error)  {
	info := OrderDepositResult{}
	params := map[string]interface{}{
		"userId":userId,//被下单人
		"orderType":orderType,
		"buyNum":buyNum,
		"tradeCode":"02",
		"busiType":"06",
		"clientIp":"127.0.0.1",
		"buyType":2,
		"couponSn":strconv.FormatInt(no,10),
		"mercId":"888000000000004",
		"platform":"ZBLMALL",
		"sysCnl":"H5",
		"timestamp":strconv.FormatInt(time.Now().Unix(),10),
	}
	res,err  := httpApiPost(viper.GetString("java.java_order_deposit_submit"),params,"v1",userToken)
	json.Unmarshal(res, &info)
	return info,err
}
package java

import (
	"encoding/json"
	"github.com/spf13/viper"
)

// 订单信息
type OrderResult struct {
	Msg					 string     `json:"msg"`
	Code				 int64		`json:"code"`
	Data 				 struct{
		OrderInfo  struct{
			Id 		 	 		 int64 		 `json:"id"`  //订单id
			OrderNo 		 	 string 	 `json:"orderNo"`  //订单编号
			GoodsTotalPrice 	 string 	 `json:"goodsTotalPrice"`  //商品总金额
			TotalPrice 			 string 	 `json:"totalPrice"`  //订单总金额
			Status		    	 string    	 `json:"status"` //订单状态 未支付/待支付:W 已完成:S 已取消:C 待发货:WS 待收货:WD
		}
	}		`json:"data"`
}
func GetOrderInfo(orderNo string)(OrderResult,error)  {
	info := OrderResult{}
	params := map[string]interface{}{
		"orderNo":orderNo,
	}
	res,err  := httpApiPost(viper.GetString("java.java_user_order_url"),params,"v1","")
	json.Unmarshal(res, &info)
	return info,err
}


// 订单购买数量统计（无token）
type OrderBuyTypeCount struct {
	Msg					 string     `json:"msg"`
	Code				 int64		`json:"code"`
	Data 				 struct{
		BuyTypeCount 		  int64 		 `json:"buyTypeCount"`  //数量
	}		`json:"data"`
}
func GetOrderBuyTypeCount(userId int64)(OrderBuyTypeCount,error)  {
	info := OrderBuyTypeCount{}
	params := map[string]interface{}{
		"userId":userId,
	}
	res,err  := httpApiPost(viper.GetString("java.java_order_buy_type_count"),params,"v1","")
	json.Unmarshal(res, &info)
	return info,err
}

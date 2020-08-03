package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mlive/service/admin"
	"mlive/service/java"
	"mlive/util"
	"strconv"
	"time"
)



// 运营后台修改云仓
func SaveCloudNumber(c *gin.Context)  {

	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}
	adminInfo			  ,_:= admin.CheckAdminLogin(adminToken)
	var adminUserId int64 = int64(adminInfo.Id)
	if  adminUserId <= 0 {
		util.Fail(c, 11111, "get admin userId error")
		return
	}

	checkParams   := []interface{}{"userId","type","number"}
	params,err 	  := util.GetRawData(c.Request,checkParams)
	if err != nil {
		util.Fail(c, 11111, fmt.Sprintf("%s",err))
		return
	}

	userParams := util.ChectIntFloat(params["userId"])
	typeParams := util.ChectIntFloat(params["type"])
	numberParams := util.ChectIntFloat(params["number"])
	if !userParams || !typeParams || !numberParams {
		util.Fail(c, 11111, "parameter fail")
		return
	}
	userId := int64(params["userId"].(float64))
	// type 1:新增,2 减少
	types := int64(params["type"].(float64))
	number := int64(params["number"].(float64))

	var remarks  string = ""
	switch params["remarks"].(type) {
	case float64:
		remarks =  strconv.FormatInt(int64(params["remarks"].(float64)),10)
	case string:
		remarks = params["remarks"].(string)
	}


	if types <= 0 || types >=3 {
		util.Fail(c, 80000, "type 类型错误")
		return
	}
	info,_ := user.GetUserInfo(userId)
	if info.Id <= 0 {
		util.Fail(c, 80000, "用户不存在")
		return
	}
	if info.Level == 0 {
		util.Fail(c, 80000, "普通用户不可以修改云仓")
		return
	}
	if number == 0 {
		util.Fail(c, 80000, "请输入数量...")
		return
	}
	cloudNum,_ := cloud.GetCloudStorageSurplusNum(userId)
	if types == 2 && number > cloudNum {
		util.Fail(c, 80000, "减少的云仓数量不足...")
		return
	}
	if  types == 2 {
		number = -number
	}
	no ,_ := cloud.CreateCloudLog(userId,types,"",0,-1,number,3,1,adminUserId,remarks)
	if no > 0 {
		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "处理成功",
			"data": map[string]interface{}{
				"no": no,
			},
		})
	}else{
		util.Fail(c, 80000, "处理失败...")
		return
	}

}

type CloudPayLog struct {
	PayNo		string 		 `json:"payNo"` //支付流水
	UserId    	int64		 `json:"userId"`
	OrderNo		string		 `json:"orderNo"`
	Mobile		string 		 `json:"mobile"`
	Nickname	string	 	 `json:"nickname"`
	Price		float64 	 `json:"price"`//金额
	OrderStatus string 		 `json:"orderStatus"`	//支付状态
	BankName	string 		 `json:"bankName"`//付款方式
	OpenId		string 		 `json:"openId"` //付款账号,银行用户标识
	PayDate		string 		 `json:"payDate"` //支付时间-支付日期 yyyy-MM-dd
	PayTime		string 		 `json:"payTime"` //支付时间-支付时间 HH:mm:ss
}

//云仓交易记录-展示 充值vip、店长、总监、合伙人、联仓的记录 (type=1&order_no!=""&cloud_type>=1)
func GetCloudPayLog(c *gin.Context)  {

	userId ,_ 		:= strconv.ParseInt(c.DefaultQuery("userId","0"),10,0)
	mobile	  		:= c.DefaultQuery("mobile","")
	row ,_    		:= strconv.Atoi(c.DefaultQuery("row","10"))
	page ,_   		:= strconv.Atoi(c.DefaultQuery("page","1"))
	code	  		:= 10000
	msg  	  		:= "success"

	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}

	mapBind := map[string]interface{}{
	}
	if userId > 0 {
		mapBind["userId"] = userId
	}
	if mobile != ""{
		mapBind["mobile"] = mobile
		mobInfo,_ := user.GetUserMobileInfo(mobile)
		if mobInfo.Id > 0 {
			mapBind["userId"] = mobInfo.Id
		}else{
			mapBind["userId"] = -1
		}
	}

	list,err := cloud.GetAdminPayCloudStoragList(page,row,mapBind)
	if err != nil {
		util.Fail(c, 80000, "get user list fail:" + fmt.Sprintf("%v", err))
		return
	}
	var data []CloudPayLog
	count,_ := cloud.GetAdminPayCloudStoragCount(mapBind)
	if len(list) > 0 {
		for _,v := range list {
			var (
				nickname 	  = ""
				mobile  	  = ""
			)
			userInfo,_ := user.GetUserInfo(v.UserId)
			if userInfo.Id > 0 {
				nickname = userInfo.Nickname
				mobile = userInfo.Mobile
			}
			var(
				payNo = ""
				price float64 = 0
				orderStatus = ""
				payDate = ""
				payTime = ""
				bankName  = ""
				openId  = ""
			)
			if v.OrderNo != "" {
				depositInfo,_ := java.GetDepositInfo(v.OrderNo)
				if depositInfo.Code == 10000 {
					payNo 		= depositInfo.Data.PayNo
					price 		= depositInfo.Data.Price
					orderStatus = depositInfo.Data.OrderStatus
					payDate 	= depositInfo.Data.PayDate
					payTime 	= depositInfo.Data.PayTime
					bankName 	= depositInfo.Data.BankName
					openId 		= depositInfo.Data.OpenId
				}
			}
			data = append(data,CloudPayLog{
				PayNo:payNo,
				UserId:v.UserId,
				Nickname:nickname,
				OrderNo:v.OrderNo,
				Mobile:mobile,
				Price:price,
				OrderStatus:orderStatus,
				BankName:bankName,
				OpenId:openId,
				PayDate:payDate,
				PayTime:payTime,
			})
		}
	}else{
		data = []CloudPayLog{}
	}
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": map[string]interface{}{
			"items": data,
			"total": count,
		},
	})

}

type CloudAllLog struct {
	UserId				int64 		`json:"userId"`
	Nickname			string	 	`json:"nickname"`
	OrderNo			    string		`json:"orderNo"`
	Level				int64 		`json:"level"`
	Type				int64 		`json:"type"`
	Channel				int64 		`json:"channel"`
	Number				int64 		`json:"number"`
	UseUserId			int64 		`json:"useUserId"`
	UseNickname			string 		`json:"useNickname"`
	CreateTime			string 		`json:"createTime"`
	UpdateTime			string 		`json:"updateTime"`
}


//云仓仓库记录-运营后台操作云仓记录、A用户赠送B用户记录、购买充值记录 (出库:减少,入库:新增)
func GetCloudAllLog(c *gin.Context)  {
	userId ,_ 		:= strconv.ParseInt(c.DefaultQuery("userId","0"),10,0)
	useUserId ,_ 	:= strconv.ParseInt(c.DefaultQuery("useUserId","0"),10,0)
	types ,_ 		:= strconv.ParseInt(c.DefaultQuery("type","0"),10,0)
	orderNo	  		:= c.DefaultQuery("orderNo","")
	row ,_    		:= strconv.Atoi(c.DefaultQuery("row","10"))
	page ,_   		:= strconv.Atoi(c.DefaultQuery("page","1"))
	code	  		:= 10000
	msg  	  		:= "success"

	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}

	mapBind := map[string]interface{}{
	}
	if userId > 0 {
		mapBind["userId"] = userId
	}
	if useUserId > 0 {
		mapBind["useUserId"] = useUserId
	}
	if types > 0 {
		mapBind["type"] = types
	}
	if orderNo != ""{
		mapBind["orderNo"] = orderNo
	}

	list,err := cloud.GetAdminAllCloudStoragList(page,row,mapBind)
	if err != nil {
		util.Fail(c, 80000, "get user list fail:" + fmt.Sprintf("%v", err))
		return
	}
	var data []CloudAllLog
	count,_ := cloud.GetAdminAllCloudStoragCount(mapBind)
	if len(list) > 0 {
		for _,v := range list {
			var (
				nickname 	  = ""
				level 	int64 = 0
				useNickname   = ""
			)
			userInfo,_ := user.GetUserInfo(v.UserId)
			if userInfo.Id > 0 {
				nickname = userInfo.Nickname
				level 	 = userInfo.Level
			}
			UseUserInfo,_ := user.GetUserInfo(v.UseUserId)
			if UseUserInfo.Id > 0 {
				useNickname = UseUserInfo.Nickname
			}
			if v.Number < 0 {
				v.Number = -v.Number
			}
			data = append(data,CloudAllLog{
				UserId:v.UserId,
				Nickname:nickname,
				OrderNo:v.OrderNo,
				Level:level,
				Type:v.Type,
				Channel:v.Channel,
				Number:v.Number,
				UseUserId:v.UseUserId,
				UseNickname:useNickname,
				CreateTime:v.CreateTime.Format("2006-01-02 15:04:05"),
				UpdateTime:v.UpdateTime.Format("2006-01-02 15:04:05"),
			})
		}
	}else{
		data = []CloudAllLog{}
	}
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": map[string]interface{}{
			"items": data,
			"total": count,
		},
	})

}


//运营后台退款-订单 java 调用
func Refund(c *gin.Context)  {
	code	      := 10000
	msg  	      := "success"

	checkParams   := []interface{}{"orderNo","adminUserId"}
	params,err 	  := util.GetRawData(c.Request,checkParams)
	if err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c,11111,msg)
		return
	}
	adminUserIdParams 	 := util.ChectIntFloat(params["adminUserId"])
	if   !adminUserIdParams {
		util.Fail(c, 11111, "arameter fail")
		return
	}
	adminUserId  := int64(params["adminUserId"].(float64))
	var orderNo string = ""
	switch params["orderNo"].(type) {
	case string:
		orderNo,_ = params["orderNo"].(string)
	}
	if orderNo == "" {
		util.Fail(c,80000,"orderNo type err or null ")
		return
	}
	orderInfo,_ := cloud.GetCloudStoragOrderNoInfo(orderNo,0)
	if orderInfo.Id <= 0 {
		util.Fail(c,80000,"no fail")
		return
	}
	if orderInfo.Status != 1 {
		util.Fail(c,80000,"status != 1 ")
		return
	}
	// 订单回收退款
	res,_ := cloud.SaveCloudLogByOrderNo(orderNo, map[string]interface{}{
		"status":-1,
		"update_time":time.Now(),
		"admin_user_id":adminUserId,
	})
	if !res {
		util.Fail(c,80000," fail")
		return
	}
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data":"",
	})
}
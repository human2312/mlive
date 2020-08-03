package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	cloudDao "mlive/service/cloud/dao"
	"mlive/service/java"
	"mlive/util"
	"strconv"
	"time"
)


var (
	cloud cloudDao.CloudStorage
)

// 提交商品云仓
func SubmitCloud(c *gin.Context)  {



	// type 0 兑换3个-->不升级 ,1 兑换21个-->升级vip,
	// 		2 兑换147个-->店长, 3 兑换588个-->总监, 4 兑换1764个-->合伙人
	checkParams   	 := []interface{}{"areaCode","mobile","orderNo"}
	code,msg,params	  := toCheckCloud(c,checkParams,-1)
	if code != 10000 {
		util.Fail(c, code, msg)
		return
	}

	userToken := c.Request.Header.Get("X-MPMall-Token")
	userId ,err:= java.Token2Id(userToken)
	if userId <= 0 {
		util.Fail(c,80000," user token fail:"+err.Msg)
		return
	}

	var types int64 = 0

	var areaCode  string = ""
	switch params["areaCode"].(type) {
	case float64:
		areaCode =  strconv.FormatInt(int64(params["areaCode"].(float64)),10)
	case string:
		areaCode = params["areaCode"].(string)
	}
	if areaCode == "" {
		util.Fail(c,80000,"areaCode number acquisition failed")
		return
	}

	var mobile  string = ""
	switch params["mobile"].(type) {
	case float64:
		mobile =  strconv.FormatInt(int64(params["mobile"].(float64)),10)
	case string:
		mobile = params["mobile"].(string)
	}

	var orderNo string = ""
	switch params["orderNo"].(type) {
	case float64:
		orderNo =  strconv.FormatInt(int64(params["orderNo"].(float64)),10)
	case string:
		orderNo = params["orderNo"].(string)
	}
	if orderNo == "" {
		util.Fail(c,80000,"orderNo to null")
		return
	}

	info,_ := user.GetUserInfo(userId)
	if info.Id <= 0 {
		util.Fail(c,80000," UserId  failed")
		return
	}

	inviteOutInfo,_ := user.GetUserAreaMobileInfo(areaCode,mobile)
	if inviteOutInfo.Id <= 0 {
		util.Fail(c,80000," 该手机号码尚未注册。请先发送邀请海报给你的好友，提醒好友注册")
		return
	}
	var number int64 = -3
	no,_ := cloud.CreateCloudLog(userId,2,orderNo,inviteOutInfo.Id,types,number,1,0,0,"")
	if no <= 0 {
		util.Fail(c,80000,"save cloud no fail...")
		return
	}
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data":map[string]interface{}{
			"no":no,
		},
	})

}

// 处理云仓
func HandleCloud(c *gin.Context)  {


	checkParams   	  := []interface{}{"type","areaCode","mobile"}
	code,msg,params	  := toCheckCloud(c,checkParams,1)


	userToken := c.Request.Header.Get("X-MPMall-Token")
	userId ,err := java.Token2Id(userToken)
	if userId <= 0 {
		util.Fail(c,80000," user token fail:"+err.Msg)
		return
	}

	//userId  	:= int64(params["userId"].(float64))
	types  		:= int64(params["type"].(float64))

	if types >=1 && types <= 4 {
		if code != 10000 {
			util.Fail(c,code,msg)
			return
		}

		var areaCode  string = ""
		switch params["areaCode"].(type) {
		case float64:
			areaCode =  strconv.FormatInt(int64(params["areaCode"].(float64)),10)
		case string:
			areaCode = params["areaCode"].(string)
		}
		if areaCode == "" {
			util.Fail(c,80000,"areaCode number acquisition failed")
			return
		}

		var mobile  string = ""
		switch params["mobile"].(type) {
		case float64:
			mobile =  strconv.FormatInt(int64(params["mobile"].(float64)),10)
		case string:
			mobile = params["mobile"].(string)
		}
		inviteOutInfo,_ := user.GetUserAreaMobileInfo(areaCode,mobile)



		// 进行升级
		var number int64 = 0
		if types == 1 {
			number = -21
		}else if types == 2 {
			number = -147
		}else if types == 3 {
			number = -588
		}else if types == 4 { // 1764 修改 5294
			number = -5294
		}
		if number == 0 {
			util.Fail(c,80000," number For 0 ")
			return
		}
		//扣云仓上级
		no,_ := cloud.CreateCloudLog(userId,2,"",inviteOutInfo.Id,types,number,1,1,0,"")
		if no <= 0 {
			util.Fail(c,80000,"save cloud no fail...")
			return
		}
		// 通知java 创建订单
		var buyNum int64 = number
		if number < 0 {
			buyNum = -buyNum
		}
		orderInfo,_ := java.SubmitOrderDeposit(no,types,inviteOutInfo.Id,buyNum,userToken)
		if orderInfo.Data.OrderNo == "" {
			util.Fail(c,80000,"java create order no error ")
			return
		}
		// 创单更新-更新订单号
		res,_ := cloud.SaveCloudLogByOrder(0, no, map[string]interface{}{
			"order_no":    orderInfo.Data.OrderNo,
			"update_time": time.Now(),
		})
		if !res {
			util.Fail(c,80000," no and order_no update error ")
			return
		}
		// 下级升级-并且赠送指标
		res,msg := user.SaveLevel(inviteOutInfo.Id,0,1,types,0,"",orderInfo.Data.OrderNo)
		if !res {
			util.Fail(c,80000,"upgrade failed."+msg)
			return
		}else{
			// 下级增加云仓
			noInvite,_ := cloud.CreateCloudLog(inviteOutInfo.Id,1,orderInfo.Data.OrderNo,userId,types,buyNum,1,1,0,"")
			// 线下通知java
			cloud.MqJavaPushUserUpdate(orderInfo.Data.OrderNo,inviteOutInfo.Id,userId,buyNum,noInvite,2)

			// 返回用户剩余的云仓数量
			cloundNum,_ := cloud.GetCloudStorageSurplusNum(userId)
			c.JSON(200, gin.H{
				"code": 10000,
				"msg":  "升级成功",
				"data":map[string]interface{}{
					"number":cloundNum,
				},
			})
		}

	}else{
		// 返回用户剩余的云仓数量
		cloundNum,_ := cloud.GetCloudStorageSurplusNum(userId)

		c.JSON(200, gin.H{
			"code": code,
			"msg":  msg,
			"data":map[string]interface{}{
				"number":cloundNum,
			},
		})
	}

}

// 封装检测
func toCheckCloud(c *gin.Context,checkParams []interface{},tc int64)(int,string,map[string]interface{})  {

	code	      := 10000
	msg  	      := "success"
	params,err 	  := util.GetRawData(c.Request,checkParams)
	if err != nil {
		msg  = fmt.Sprintf("%s",err)
		return 11111,msg,params
	}
	var types int64 = 0
	if tc == 1 {
		typeParams := util.ChectIntFloat(params["type"])
		if !typeParams {
			msg = " parameter fail"
			return 11111, msg, params
		}
		types = int64(params["type"].(float64))
		if types < 0 || types > 6 {
			return 11111, "type  error", params
		}
	}

	userToken := c.Request.Header.Get("X-MPMall-Token")
	userId ,err1 := java.Token2Id(userToken)
	if userId <= 0 {
		return 80000,"user token fail:"+err1.Msg,params
	}

	var areaCode  string = ""
	switch params["areaCode"].(type) {
	case float64:
		areaCode =  strconv.FormatInt(int64(params["areaCode"].(float64)),10)
	case string:
		areaCode = params["areaCode"].(string)
	}
	if areaCode == "" {
		return 80000,"areaCode number acquisition failed",params
	}


	var mobile  string = ""
	switch params["mobile"].(type) {
	case float64:
		mobile =  strconv.FormatInt(int64(params["mobile"].(float64)),10)
	case string:
		mobile = params["mobile"].(string)
	}
	if mobile == "" {
		return 80000,"Mobile number acquisition failed",params
	}

	info,_ := user.GetUserInfo(userId)
	if info.Id <= 0 {
		return 80000,"UserId  failed",params
	}

	inviteOutInfo,_ := user.GetUserAreaMobileInfo(areaCode,mobile)
	if inviteOutInfo.Id <= 0 {
		return 80000,"该手机号码尚未注册。请先发送邀请海报给你的好友，提醒好友注册",params
	}
	// 判断邀请下级链条
	//isSub,_ := ago.CheckIsInviteSubordinate(userId,inviteOutInfo.Id)
	//if !isSub {
	//	return 80000,"该手机号码并非你的邀请下级，请核实",params
	//}
	//判断直属邀请下级
	isSubInviTe,_ := user.GetUserCount(inviteOutInfo.Id,userId)
	if isSubInviTe <= 0 {
		return 80000,"该手机号码并非你的邀请下级，请核实",params
	}

	var number 	 int64 = 0
	//var levelName string = ""
	if types == 0 {
		if inviteOutInfo.Level >= 1 {
			return 80000,"商品只能赠送给普通用户...",params
		}
		orderInfo,err := java.GetOrderBuyTypeCount(inviteOutInfo.Id)
		if err != nil {
			return 80000,err.Error(),params
		}
		if orderInfo.Code != 10000 {
			return 80000,"java order buy msg:"+orderInfo.Msg+" error,code:"+strconv.FormatInt(orderInfo.Code,10),params
		}
		if orderInfo.Data.BuyTypeCount > 0 {
			return 80000,"该用户已存在订单记录",params
		}
		total,_	:= cloud.GetCloudStoragUseCount(inviteOutInfo.Id)
		if total > 0 {
			if inviteOutInfo.Level == 0 {
				return 80000,"该用户已存在订单记录....",params
			}
			return 80000,"商品只能给予给普通用户",params
		}
		number = 3
		//levelName = "普通用户"
	}else if types == 1 {
		number = 21
		//levelName = "vip"
	}else if types == 2 {
		number = 147
		//levelName = "店长"
	}else if types == 3 {
		number = 588
		//levelName = "总监"
	}else if types == 4 { // 1764-5294
		number = 5294
		//levelName = "合伙人"
	}else{
		return 80000," 类型参数错误...",params
	}

	var inviteOutLevelName string = "未知"
	if inviteOutInfo.Level == 0 {
		inviteOutLevelName = "普通用户"
	}else if inviteOutInfo.Level == 1 {
		inviteOutLevelName = "vip"
	}else if inviteOutInfo.Level == 2 {
		inviteOutLevelName = "店长"
	}else if inviteOutInfo.Level == 3 {
		inviteOutLevelName = "总监"
	}else if inviteOutInfo.Level == 4 {
		inviteOutLevelName = "合伙人"
	}else if inviteOutInfo.Level == 5 {
		inviteOutLevelName = "联创"
	}
	if types >=1 {
		if inviteOutInfo.Level >= types {
			return 80000," 该手机号码已经是"+inviteOutLevelName+"，无法处理",params
		}
		// 若该用户邀请的用户等级高于或者等于自己，则toast：邀请的用户必须低于你的等级
		if info.Level <= types {
			return 80000,"邀请的用户必须低于你的等级",params
		}
		//stayNum,_ := cloud.GetCloudStoragStay(inviteOutInfo.Id)
		//if stayNum > 0 {
		//	return 80000," 上级未处理旧订单支付...",params
		//}
	}
	if number <= 0 {
		return 80000,"number error",params
	}
	cloundNum,_ := cloud.GetCloudStorageSurplusNum(userId)
	if cloundNum <= 0 {
		return 80000,"数量不足",params
	}
	if cloundNum <  number {
		return 80000,"数量不足",params
	}
	return code,"success",params
}

// 云仓数量
func CloudInfo(c *gin.Context)  {

	userToken := c.Request.Header.Get("X-MPMall-Token")
	userId ,err := java.Token2Id(userToken)
	if userId <= 0 {
		util.Fail(c,80000," user token fail:"+err.Msg)
		return
	}

	num,_ 		:= cloud.GetCloudStorageSurplusNum(userId)
	var dataArr = []map[string]interface{}{}
	info,_ 		:= user.GetUserInfo(userId)
	dataMap0 := map[string]interface{}{"type":0,"name":"赠送商品"}
	dataMap1 := map[string]interface{}{"type":1,"name":"VIP套餐"}
	dataMap2 := map[string]interface{}{"type":2,"name":"店长套餐"}
	dataMap3 := map[string]interface{}{"type":3,"name":"总监套餐"}
	dataMap4 := map[string]interface{}{"type":4,"name":"合伙人套餐"}

	if info.Level == 1{
		dataArr = append(dataArr,dataMap0)
	}
	if info.Level == 2{
		dataArr = append(dataArr,dataMap0,dataMap1)
	}
	if info.Level == 3{
		dataArr = append(dataArr,dataMap0,dataMap1,dataMap2)
	}
	if info.Level == 4{
		dataArr = append(dataArr,dataMap0,dataMap1,dataMap2,dataMap3)
	}
	if info.Level == 5{
		dataArr = append(dataArr,dataMap0,dataMap1,dataMap2,dataMap3,dataMap4)
	}

	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data":map[string]interface{}{
			"num":num,
			"column":dataArr,
		},
	})
}

//回收商品
func Recovery(c *gin.Context)  {

	code	      := 10000
	msg  	      := "success"
	checkParams   := []interface{}{"no"}
	params,err 	  := util.GetRawData(c.Request,checkParams)
	if err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c,11111,msg)
		return
	}
	userToken := c.Request.Header.Get("X-MPMall-Token")
	userId ,err1 := java.Token2Id(userToken)
	if userId <= 0 {
		util.Fail(c,80000," user token fail:"+err1.Msg)
		return
	}
	var cloudNo int64 = 0
	switch params["no"].(type) {
	case float64:
		cloudNo =  int64(params["no"].(float64))
	case string:
		cloudNo,_ = strconv.ParseInt(params["no"].(string),10,64)
	}
	if cloudNo <= 0 {
		util.Fail(c,80000,"no to null")
		return
	}
	orderInfo,_ := cloud.GetCloudStoragOrderInfo(userId,cloudNo)
	if orderInfo.Id <= 0 {
		util.Fail(c,80000,"no fail")
		return
	}
	if orderInfo.CloudType != 0 {
		util.Fail(c,80000,"Cloud type  operation failed ")
		return
	}
	if orderInfo.Status < 0 {
		util.Fail(c,80000,"no expired")
		return
	}
	if orderInfo.Status >= 1 {
		util.Fail(c,80000,"no  finish")
		return
	}
	// 回收
	res,_ := cloud.SaveCloudLogByOrder(userId,cloudNo, map[string]interface{}{
		"status":-1,
		"update_time":time.Now(),
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

// 返回数据
type CloudData struct {
	UserId 			int64 		`json:"userId"`
	HeadImgUrl 		string 		`json:"headImgUrl"`
	Nickname 		string 		`json:"nickname"`
	CreateTime 		string 		`json:"createTime"`
	Number			int64		`json:"number"`
}

// 获取列表
func CloudList(c *gin.Context)  {

	//5hwipp8nhb0moh4qs13p119f50tulp6o1
	userToken := c.Request.Header.Get("X-MPMall-Token")
	userId ,err1 := java.Token2Id(userToken)
	if userId <= 0 {
		util.Fail(c,80000," user token fail:"+err1.Msg)
		return
	}
	checkParams   	 := []interface{}{"limit","page"}
	params,err 	  := util.GetRawData(c.Request,checkParams)
	if err != nil {
		util.Fail(c,11111,fmt.Sprintf("%s",err))
		return
	}
	limitParams := util.ChectIntFloat(params["limit"])
	pageParams := util.ChectIntFloat(params["page"])
	if !limitParams  || !pageParams {
		util.Fail(c,11111," parameter fail")
	}
	limit  := int64(params["limit"].(float64))
	page  := int64(params["page"].(float64))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	list,_   := cloud.GetCloudStoragList(userId,page,limit)
	var data []CloudData
	if len(list) > 0 {
		for _, v := range list {
			inviteOutInfo,_ := user.GetUserInfo(v.UseUserId)
			if inviteOutInfo.Id > 0 {
				if v.Number < 0 {
					v.Number = -v.Number
				}
				data = append(data, CloudData{
					UserId:inviteOutInfo.Id,
					HeadImgUrl:inviteOutInfo.HeadImgUrl,
					Nickname:inviteOutInfo.Nickname,
					CreateTime:v.CreateTime.Format("2006-01-02 15:04:05"),
					Number:v.Number,
				})
			}
		}
	}
	if data == nil {
		data = []CloudData{}
	}
	count,_  := cloud.GetCloudStoragCount(userId)
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data":map[string]interface{}{
			"items": data,
			"total": count,
		},
	})
}



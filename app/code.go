package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	code "mlive/service/code/dao"
	"mlive/util"
	"mlive/service/java"
	"strconv"
)

var (
	userCode code.UserCode
)
// 名额数量
func CodeNum(c *gin.Context)  {

	// type 1 商品,2 店长,3 总监,	4 合伙人，5 (联创
	checkParams   := []interface{}{"userId","type"}
	params,err 	  := util.GetRawData(c.Request,checkParams)
	//code	      := 10000
	msg  	      := "success"
	if err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c, 11111, msg)
		return
	}

	userParams := util.ChectIntFloat(params["userId"])
	typeParams := util.ChectIntFloat(params["type"])
	if !userParams || !typeParams {
		msg = " parameter fail"
		util.Fail(c, 11111, msg)
		return
	}
	userId  := int64(params["userId"].(float64))
	ty  	:= int64(params["type"].(float64))

	log.Println("userId:",userId)
	log.Println("ty:",ty)

	if ty == 1 {

	}else if ty >=2 && ty <=5 { // 服务商升级

	}else{
		util.Fail(c,80000,"类型错误")
	}



}

// 使用码
func UseCode(c *gin.Context)  {


	// type 1 商品,2 店长,3 总监,	4 合伙人，5 (联创
	checkParams   := []interface{}{"userId","type","mobile"}
	params,err 	  := util.GetRawData(c.Request,checkParams)
	//code	      := 10000
	msg  	      := "success"
	if err != nil {
		msg  = fmt.Sprintf("%s",err)
		util.Fail(c, 11111, msg)
		return
	}

	userParams := util.ChectIntFloat(params["userId"])
	typeParams := util.ChectIntFloat(params["type"])
	if !userParams || !typeParams {
		msg = " parameter fail"
		util.Fail(c, 11111, msg)
		return
	}
	userId  := int64(params["userId"].(float64))
	ty  	:= int64(params["type"].(float64))

	var mobile  string = ""
	switch params["mobile"].(type) {
		case float64:
			mobile =  strconv.FormatInt(int64(params["mobile"].(float64)),10)
		case string:
			mobile = params["mobile"].(string)
	}
	if mobile == "" {
		util.Fail(c,80000,"手机号码获取失败")
		return
	}
	var levelName string = ""
	if ty == 1 {
		levelName = "VIP"
	}else if ty == 2 {
		levelName = "店长"
	}else if ty == 3 {
		levelName = "总监"
	}else if ty == 4 {
		levelName = "合伙人"
	}
	inviteOutInfo,_ := user.GetUserMobileInfo(mobile)
	if inviteOutInfo.Id <= 0 {
		util.Fail(c,80000," 该手机号码尚未注册。请先发送邀请海报给你的好友，提醒好友注册:"+mobile)
		return
	}
	isSub,_ := ago.CheckIsInviteSubordinate(userId,inviteOutInfo.Id)
	log.Println("isSub:",isSub)
	if !isSub {
		util.Fail(c,80000," 该手机号码并非你的邀请下级，请核实")
		return
	}
	totalQuota,_ := userCode.GetCodeLogNum(userId,ty)
	if totalQuota <= 0 {
		util.Fail(c,80000,"赠送名额不足")
		return
	}
	log.Println("userId:",userId)
	log.Println("ty:",ty)

	if ty == 0 {// 商品码 待处理
		// 只允许赠送给无订单记录、无被赠送记录的普通用户，判断被赠送用户是否为普通用户，则toast：商品只能赠送给普通用户
		//若该用户有订单、被赠送记录，则toast：该用户已存在订单记录
		//GetOrderBuyTypeCount
		orderInfo,err := java.GetOrderBuyTypeCount(inviteOutInfo.Id)
		if err != nil {
			util.Fail(c,80000,err.Error())
			return
		}
		if orderInfo.Data.BuyTypeCount > 0 {
			util.Fail(c,80000,"该用户已存在订单记录")
			return
		}
		totalScore,_ := userCode.GetCodeLogCount(inviteOutInfo.Id)
		if totalScore >0 {
			util.Fail(c,80000,"商品只能赠送给普通用户....")
			return
		}
		no,_ := userCode.CreateCodeLog(userId,2,ty,inviteOutInfo.Id,-1,1,0,0)
		if no <= 0 {
			util.Fail(c,80000,"扣除名额失败")
			return
		}
		mapQuota := map[string]interface{}{
			"goods":gorm.Expr("goods-?",1),
		}
		res,_ := userCode.SaveCodeSummaryColumn(userId,mapQuota)
		if !res {
			util.Fail(c,80000,"更新汇总名额失败")
			return
		}

		c.JSON(200, gin.H{
			"code": 10000,
			"msg":  "ok",
			"data":map[string]interface{}{
				"no":no,
			},
		})
		return

	}else if ty >=1 && ty <=4 { // 服务商升级
		if ty == inviteOutInfo.Level {
			util.Fail(c,80000," 该手机号码已经是"+levelName+"，无法处理")
			return
		}
		// 扣除赠送的名额
		no,_ := userCode.CreateCodeLog(userId,2,ty,inviteOutInfo.Id,-1,1,1,0)
		if no <= 0 {
			util.Fail(c,80000,"扣除名额失败")
			return
		}
		mapQuota := map[string]interface{}{}
		if ty == 1 {
			mapQuota["vip"] =  gorm.Expr("vip-?",1)
		}else if ty == 2 {
			mapQuota["shop_owner"] =  gorm.Expr("shop_owner-?",1)
		}else if ty == 3 {
			mapQuota["chief_inspector"] =  gorm.Expr("chief_inspector-?",1)
		}else if ty == 4 {
			mapQuota["partner"] = gorm.Expr("partner-?",1)
		}
		if len(mapQuota) >0 {
			res,_ := userCode.SaveCodeSummaryColumn(userId,mapQuota)
			if !res {
				util.Fail(c,80000,"更新汇总名额失败")
				return
			}
			res1,msg := user.SaveLevel(inviteOutInfo.Id,0,1,ty,0 ,"","")
			if !res1 {
				util.Fail(c,80000,msg)
				return
			}

			c.JSON(200, gin.H{
				"code": 10000,
				"msg":  "ok",
				"data":map[string]interface{}{
					"no":no,
				},
			})
			return

		}else{
			util.Fail(c,80000,"减少汇总名额失败")
			return
		}

	}else{
		util.Fail(c,80000,"操作类型错误")
		return
	}
}

// 检查名额是否合法
func checkCode(c *gin.Context)  {

}


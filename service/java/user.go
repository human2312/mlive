package java

import (
	"encoding/json"
	"github.com/spf13/viper"
	"strconv"
)

type JavaUserInfoResult struct {
	UserId 		 	 int64 		`json:"id"`  //用户id
	InviteUserId 	 int64 		`json:"inviteUserId"`  //邀请上级id
	UserName 	 	 string 	`json:"username"`//用户登录名
	HeadImgUrl 	 	 string 	`json:"headImgUrl"`//头像url地址
	Nickname	 	 string 	`json:"nickname"`//昵称
	Name	 	 	 string 	`json:"name"`//姓名
	Gender	 	  	 int64      `json:"gender"`//性别 0:未知 1:男 2:女
	PtLevel	 	  	 int64      `json:"ptLevel"`//用户等级(0 粉丝、1 主播、2 服务商、3 总监、4 合伙人、5 基地)
	InviteCode	 	 string     `json:"inviteCode"`//邀请码
	AreaCode	 	 string     `json:"areaCode"`//区号
	Mobile	 	 	 string     `json:"mobile"`//手机号
	TelPhone	 	 string     `json:"telPhone"`//电话
	Email	     	 string     `json:"email"`//电子邮箱
	Status	     	 int64    	`json:"status"`//状态  1:正常 2:冻结/黑名单(封禁)
	Operator		 string		`json:"operator"`// 操作用户姓名
	AddTime		     string		`json:"addTime"`// 创建时间
	UpdateTime		 string		`json:"updateTime"`//更新信息
}

/**
* 获取用户信息
 */
func InfoById(userId int64)(JavaUserInfoResult,error)  {
	info := JavaUserInfoResult{}
	params := map[string]interface{}{
		"id":userId,
	}
	res,err  := httpApiPost(viper.GetString("java.java_user_info_url"),params,"v1","")
	json.Unmarshal(res, &info)
	return info,err
}


// 运营后台-用户账号-统计
type AdminUserAccountInfoStatistics struct {
	Msg					 string     `json:"msg"`
	Code				 int64		`json:"code"`
	Data 				 struct{
		UserId 		 	  int64 		 `json:"id"`  //用户id
		TotalBuy 		  int64 		 `json:"totalBuy"`  //累计订购量
		TotalSale 		  float64 		 `json:"totalSale"`  //累计销售
		TotalIncome 	  float64 		 `json:"totalIncome"`  //累计收益
		TotalUavaIncome   float64 		 `json:"totalUavaIncome"`  //在途收益
		WithdrawIncome    float64 		 `json:"withdrawIncome"`  //可提收益
		WithdrawAlready   float64 		 `json:"withdrawAlready"`  //已提现
	}		`json:"data"`
}

func GetAdminUserAccountInfoStatistics(userId int64,adminToken  string)(AdminUserAccountInfoStatistics ,error)  {
	info := AdminUserAccountInfoStatistics{}

	res,err  := httpAdminGet("Get",viper.GetString("java.java_admin_user_account_info")+"?id="+strconv.FormatInt(userId,10),adminToken)
	json.Unmarshal(res, &info)
	return info,err
}


type XFYLUserInfo struct {
	Msg					 string     `json:"msg"`
	Code				 int64		`json:"code"`
	Data 				 struct{
		UserId 		 	  int64 		 `json:"id"`  //用户id
		UserName 		  string 		 `json:"userName"`  //名字
		Nickname 		  string 		 `json:"nickname"`  //昵称
		PtLevel 		  int64 		 `json:"ptLevel"`  //等级 ( ptLevel，用户等级0:普通用户;1:VIP;2:县级店;3:品牌店;4:金钻;5:总裁;6:分公司)
	}		`json:"data"`
}

// 幸福引力用户信息
func GetXFYLUserInfo(areaCode string,mobile string)(XFYLUserInfo ,error)  {
	info := XFYLUserInfo{}
	params := map[string]interface{}{
		"areaCode":areaCode,
		"mobile":mobile,
	}
	res,err  := httpApiPost(viper.GetString("java.java_xfyl_user_info"),params,"v1","")
	json.Unmarshal(res, &info)
	return info,err
}

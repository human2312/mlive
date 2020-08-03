package sign

// @Time : 2020年3月8日14:17:05
// @Author : Lemyhello
// @Desc: sign

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"mlive/service/admin"
	"mlive/util"
	"sort"
	"strconv"
)

//SignCheck 验签入口
func Check(c *gin.Context) {
	xSign := c.Request.Header.Get("X-MPMALL-Sign")
	appVer := c.Request.Header.Get("X-MPMALL-SignVer")
	//获取Body
	var jsonBody map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request.Body)
	// 把刚刚读出来的再写进去
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	json.Unmarshal(body, &jsonBody)
	if jsonBody["mercId"] == nil || jsonBody["platform"] == nil || jsonBody["sysCnl"] == nil || jsonBody["timestamp"] == nil {
		util.Fail(c, 11111, "缺少参数或参数为空_public")
		return
	}
	if xSign == "" || appVer == "" {
		util.Fail(c, 11111, "缺少参数或参数为空_raw")
		return
	}
	plain := ""
	//进行排序
	var keys []string
	for k := range jsonBody {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if jsonBody[k] == nil {
			continue
		}
		switch jsonBody[k].(type) {
		case int:
			plain += k + "=" + strconv.Itoa((jsonBody[k]).(int)) + "&"
		case int64:
			plain += k + "=" + strconv.Itoa((jsonBody[k]).(int)) + "&"
		case string:
			plain += k + "=" + jsonBody[k].(string) + "&"
		case float64:
			plain += k + "=" + (strconv.FormatFloat(jsonBody[k].(float64), 'f', -1, 64)) + "&"
		case float32:
			plain += k + "=" + (strconv.FormatFloat(jsonBody[k].(float64), 'f', -1, 32)) + "&"
		default:
			continue
		}
	}
	//签名
	plain += "key=" + viper.GetString("sign.app_sign_ver_key_"+appVer)
	//fmt.Println(plain)
	checkSign := util.Md5V(plain)
	checkSign = util.EncodeBase64(checkSign)
	//鉴权
	if checkSign != xSign {
		util.Fail(c, 99999, "鉴权失败")
		return
	}
	return
}

//SignSet 生成签名
func Set(params map[string]interface{}, appVer string) string {
	if len(params) <= 0 {
		return ""
	}
	plain := ""
	//进行排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if params[k] == nil {
			continue
		}
		switch params[k].(type) {
		case int:
			plain += k + "=" + strconv.Itoa((params[k]).(int)) + "&"
		case int64:
			plain += k + "=" +strconv.FormatInt(params[k].(int64),10) + "&"
		case string:
			plain += k + "=" + params[k].(string) + "&"
		case float64:
			plain += k + "=" + (strconv.FormatFloat(params[k].(float64), 'f', -1, 64)) + "&"
		case float32:
			plain += k + "=" + (strconv.FormatFloat(params[k].(float64), 'f', -1, 32)) + "&"
		default:
			continue
		}
	}
	//签名
	plain += "key=" + viper.GetString("sign.app_sign_ver_key_"+appVer)
	checkSign := util.Md5V(plain)
	checkSign = util.EncodeBase64(checkSign)
	return checkSign
}

func Admin(c *gin.Context)  {
	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}
	adminInfo,err:= admin.CheckAdminLogin(adminToken)
	var adminUserId int64 = int64(adminInfo.Id)
	if  adminUserId <= 0 {
		util.Fail(c, err.Code, err.Msg)
		return
	}
}
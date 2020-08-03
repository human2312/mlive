package java

import (
	"encoding/json"
	"github.com/spf13/viper"
	"fmt"
	error2 "mlive/library/error"
)

/**
 * @Author: Lemyhello
 * @Description:
 * @File:  token
 * @Version: X.X.X
 * @Date: 2020/3/25 上午11:07
 */

//Token2IdStruct token转uid对象
type Token2IdStruct struct {
	Code int `json:"code"`
	Data struct {
		ID int64 `json:"id"`
	} `json:"data"`
	Msg string `json:"msg"`
}

//Token2Id token转uid userId = 0为无效或者未登录
func Token2Id(XMPMallToken string) (userId int64,err *error2.Err) {
	data := Token2IdStruct{}
	params := map[string]interface{}{
		"token":XMPMallToken,
	}
	res,err1  := httpApiPost(viper.GetString("java.java_api_url_token2Id"),params,"v1","")
	if err1 != nil {
		return 0,error2.New(80000,fmt.Sprintf("%v",err1))
	}
	json.Unmarshal(res, &data)
	if data.Code == 10000 && data.Data.ID > 0 {
		err = nil
	} else {
		err = error2.New(data.Code,data.Msg)
	}
	return data.Data.ID,err
}

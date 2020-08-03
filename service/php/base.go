package php

import (
	"encoding/json"
	"github.com/kirinlabs/HttpRequest"
	"github.com/spf13/viper"
	"log"
	"mlive/library/sign"
	"strconv"
	"mlive/library/logger"
	"errors"
)

type InvestorsDetail  struct {
	Msg					 string     `json:"msg"`
	Code				 int64		`json:"code"`
	Data 				 struct{
		Items 	[]InvestorsDetailItems `json:"items"`
	}		`json:"data"`
}
type InvestorsDetailItems struct {
	Id       string   `json:"id"`
	Plain    string `json:"plain"` //普通用户价格
	Plus     string `json:"plus"` // vip价格
	Train    string `json:"train"` //店长价格
	Serv     string `json:"serv"`//总监价格
	Partner  string `json:"partner"`//合伙人价格
	Director string `json:"director"`//联创价格
}

func GetInvestorsList()(InvestorsDetail,error)  {
	info := InvestorsDetail{}
	params := map[string]interface{}{}
	res,err  := httpApiPost(viper.GetString("php.php_investGoodsDetail_url"),params,"v1")
	log.Println("res:",string(res))
	json.Unmarshal(res, &info)
	return info,err
}


func httpApiPost(url string,postData map[string]interface{},ver string) (result []byte,err error) {
	// 组装sign
	sign := sign.Set(postData,ver)
	req := HttpRequest.NewRequest()
	req.SetTimeout(5) // 超时5秒
	mapHeader := map[string]string{
		"Content-Type":"application/json",
		"X-MPMALL-SignVer":ver,
		"X-MPMALL-Sign":sign,
	}
	req.SetHeaders(mapHeader)

	resp,err := req.Post(url,postData)
	if err != nil {
		logger.Eprintln(err)
	}
	body, err := resp.Body()
	//body,err := resp.Body()
	if err != nil {
		logger.Eprintln(err)
	}
	if resp.StatusCode() != 200 {
		statusStr := strconv.Itoa(resp.StatusCode())
		logger.Eprintln("接口:",url,"通信失败,返回码:",statusStr,mapHeader,string(body))
		return body,errors.New("java url error:"+url+",code:"+statusStr)
	}
	return body,nil
}

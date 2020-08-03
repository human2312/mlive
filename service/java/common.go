package java

import (
	"errors"
	"github.com/kirinlabs/HttpRequest"
	"mlive/library/logger"
	"mlive/library/sign"
	"strconv"
)

func httpAdminGet(method string,url string,adminToken string) ([]byte,error) {

	req := HttpRequest.NewRequest()
	req.SetTimeout(5) // 超时5秒
	mapHeader := map[string]string{
		"Content-Type":"application/json",
		"Admin-Token":adminToken,
	}
	req.SetHeaders(mapHeader)
	resp,err := req.Get(url)
	if err != nil {
		logger.Eprintln(err)
		logger.Eprintln("接口:",url,"通信失败resp err,返回码",mapHeader)
		return nil,err
	}
	statusStr := strconv.Itoa(resp.StatusCode())
	if resp.StatusCode() != 200 {
		logger.Eprintln("接口:",url,"通信失败,返回码",statusStr,mapHeader)
		return nil,errors.New("java error url :"+url+",code:"+statusStr)
	}

	body, err := resp.Body()
	//body,err := resp.Body()
	if err != nil {
		logger.Eprintln(err)
		logger.Eprintln("接口:",url,"通信失败body err,返回码",statusStr,mapHeader,string(body))
		return nil,err
	}
	logger.Iprintln("接口:",url,"通信,返回码",statusStr,mapHeader,string(body))
	return body,nil
}


func httpApiPost(url string,postData map[string]interface{},ver string,userToken string) (result []byte,err error) {
	// 组装sign
	sign := sign.Set(postData,ver)
	req := HttpRequest.NewRequest()
	req.SetTimeout(5) // 超时5秒
	mapHeader := map[string]string{
		"Content-Type":"application/json",
		"X-MPMALL-SignVer":ver,
		"X-MPMALL-Sign":sign,
	}
	if userToken != "" {
		mapHeader["X-MPMall-Token"] = userToken
	}

	req.SetHeaders(mapHeader)

	resp,err := req.Post(url,postData)
	if err != nil {
		logger.Eprintln(err)
		logger.Eprintln("接口:",url,"通信失败 resp err,返回码:",mapHeader,postData)
		return nil,err
	}
	statusStr := strconv.Itoa(resp.StatusCode())
	if resp.StatusCode() != 200 {
		logger.Eprintln("接口:",url,"通信失败,返回码:",statusStr,mapHeader,postData)
		return nil,errors.New("java url error:"+url+",code:"+statusStr)
	}
	body, err := resp.Body()
	if err != nil {
		logger.Eprintln(err)
		logger.Eprintln("接口:",url,"通信失败 body  err,返回码:",statusStr,mapHeader,postData,string(body))
		return nil,err
	}
	logger.Iprintln("接口:",url,"通信,返回码",statusStr,mapHeader,postData,string(body))
	return body,nil
}

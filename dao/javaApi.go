package dao

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	error2 "mlive/library/error"
	"mlive/library/logger"
	"net/http"
	"strconv"
	"strings"
)

type AdminGenerated struct {
	Msg string `json:"msg"`
	Code int `json:"code"`
	Data struct {
		ID int `json:"id"`
		Username string `json:"username"`
	} `json:"data"`
}

type AdminInfoGenerated struct {
	Msg string `json:"msg"`
	Code int `json:"code"`
	Data struct {
		ID int `json:"id"`
		Username string `json:"username"`
		LastLoginIP string `json:"lastLoginIp"`
		Avatar string `json:"avatar"`
		AddTime string `json:"addTime"`
		UpdateTime string `json:"updateTime"`
		Deleted int `json:"deleted"`
		Version int `json:"version"`
		RoleID int `json:"roleId"`
	} `json:"data"`
}

func GetAdminInfo(adminToken string) (AdminGenerated,*error2.Err) {
	admininfo := AdminGenerated{}
	res ,err:= HttpPost("POST",viper.GetString("java.java_admininfo_url"),adminToken,"")
	json.Unmarshal(res, &admininfo)
	return admininfo,err
}

func GetAdminInfoById(adminToken string,id int) (AdminInfoGenerated,*error2.Err)  {
	admininfo := AdminInfoGenerated{}
	res ,err:= HttpPost("GET",viper.GetString("java.java_adminread_url")+"?AdminId="+strconv.Itoa(id),adminToken,"")
	json.Unmarshal(res, &admininfo)
	return admininfo,err
}

func HttpPost(method string,url string,adminToken string,params string) (result []byte,err *error2.Err) {
	client := &http.Client{}
	req, neterr := http.NewRequest(method, url, strings.NewReader(params))
	if neterr != nil {
		logger.Eprintln(neterr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Admin-Token", adminToken)
	resp, resperr := client.Do(req)
	defer resp.Body.Close()
	if resperr != nil {
		logger.Eprintln(resperr)
	}
	if resp.StatusCode != 200 {
		logger.Eprintln("接口:",url,"通信失败,返回码",resp.StatusCode,resp.Header,resp.Request)
		err = error2.New(resp.StatusCode,"网络通信接口出错")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body,err
}

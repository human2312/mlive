package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func EncodeBase64(plain string) (cipher string) {
	const key string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var (
		length        int = len(plain)
		loop_time     int = length / 3
		left          int = length % 3
		bytes_catnate int32
		index         int
	)
	for i := 0; i < loop_time; i++ {
		bytes_catnate = int32(plain[index])<<16 + int32(plain[index+1])<<8 + int32(plain[index+2])
		index += 3
		cipher += string(key[(bytes_catnate>>18)&0x3F])
		cipher += string(key[(bytes_catnate>>12)&0x3F])
		cipher += string(key[(bytes_catnate>>6)&0x3F])
		cipher += string(key[bytes_catnate&0x3F])
	}
	if left == 1 {
		bytes_catnate = int32(plain[index]) << 4
		cipher += string(key[(bytes_catnate>>6)&0x3F])
		cipher += string(key[bytes_catnate&0x3F])
		cipher += "=="
	}
	if left == 2 {
		bytes_catnate = (int32(plain[index])<<8 + int32(plain[index+1])) << 2
		cipher += string(key[(bytes_catnate>>12)&0x3F])
		cipher += string(key[(bytes_catnate>>6)&0x3F])
		cipher += string(key[bytes_catnate&0x3F])
		cipher += "="
	}

	return cipher
}

//Fail 输出且结束
func Fail(c *gin.Context,code int,msg string)  {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": "",
	})
	c.Abort() //中断退出
}

//SendDingDing 钉钉通知
//https://oapi.dingtalk.com/robot/send?access_token=4d59a7e0e25463b7daf4bc258c78a6c360a2a5f3cb0b734de3c8a6b326fc5fb2
func SendDingDing(strArr ...interface{})(StatusCode int)  {
	sendContent := viper.GetString("gin.mode") + " 环境 "+ time.Now().Format("2006-01-02 15:04:05")
	for _, arg := range strArr {
			sendContent = sendContent + fmt.Sprintf(" %s", arg)
	}
	client := &http.Client{}
	access_token := "4d59a7e0e25463b7daf4bc258c78a6c360a2a5f3cb0b734de3c8a6b326fc5fb2"
	//access_token := "595f118068dfe09111dccbf8e2a48814da07a0507c55d4bc98c38e921a8be11f"
	url      := "https://oapi.dingtalk.com/robot/send?access_token="+access_token
	data := `{"msgtype":"text","text":{"content":"mlive 异常通知： `+strings.Replace(sendContent,"\n","",-1)+`"},"isAtAll":true}`
	request,err  := http.NewRequest("POST",url,strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	//增加header选项
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Length", strconv.Itoa(len(data)))
	//处理返回结果
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	return response.StatusCode
}
package main

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"mlive/mq"
	"mlive/util"
)

type TestGenerated struct {
	UserId 				int 		 `json:"userId"`
	Type	 			int			 `json:"type"`
	No	 				string		 `json:"no"`
	OrderType	 		string		 `json:"orderType"`
}

var (
	mqUserLevel1           = new(util.RabbitMQ)
)

//发送一条消息
func main()  {
	mq.ConfigInit()
	sendJson := TestGenerated{}
	sendJson.UserId = 5
	sendJson.Type   = 1
	sendJson.No     = "2020030800224321"
	sendJson.OrderType     = "goods"
	//转为Json
	waitSend, _ := json.Marshal(sendJson)
	//fmt.Println(string(waitSend))
	e := mqUserLevel1.Publish(viper.GetString("rabbitmqueue.monitorUserLevel"),string(waitSend))
	if e != nil {
	 	log.Println("发送失败:")
	}
	log.Println("发送成功!")
	return
}
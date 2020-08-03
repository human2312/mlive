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
}

var (
	mqUserLevel1           = new(util.RabbitMQ)
)


//发送一条消息
func main()  {
	mq.ConfigInit()
	sendJson := TestGenerated{}
	sendJson.UserId = 1
	//转为Json
	waitSend, _ := json.Marshal(sendJson)
	//fmt.Println(string(waitSend))
	e := mqUserLevel1.Publish(viper.GetString("rabbitmqueue.monitorUser"),string(waitSend))
	if e != nil {
	 	log.Println("发送失败:")
	}
	log.Println("发送成功!")
	return
}
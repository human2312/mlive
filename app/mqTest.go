package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"mlive/library/logger"
	"mlive/util"
	"strconv"
	"time"
)

var (
	RabbitMQ           = new(util.RabbitMQ)
)


type TestGenerated struct {
	UserId int `json:"userId"`
}


//发送一条消息
func SendTest(c *gin.Context)  {
	sendJson := TestGenerated{}
	var jsonBody map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request.Body)
	// 把刚刚读出来的再写进去
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	json.Unmarshal(body, &jsonBody)
	userId , _ := strconv.Atoi((jsonBody["userid"]).(string))
	sendJson.UserId = userId
	//转为Json
	waitSend, _ := json.Marshal(sendJson)
	//fmt.Println(string(waitSend))
	e := RabbitMQ.Publish("golang-test-queue2",string(waitSend))
	if e != nil {
		c.JSON(200, gin.H{
			"code": 11111,
			"msg":  "发送失败",
			"data": e,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": 1,
	})
	fmt.Println("发送成功!")
	return
}

//读到一条消息
func ReceiveTest()  {
	a := RabbitMQ.Consume("golang-test-queue2",funcTest)
	logger.Iprintln("消息回传",a)
	return
}

//funcTest 函数作为参数传递(回调)
func funcTest(message [] byte) (bool) {
	testJson := TestGenerated{}
	json.Unmarshal(message, &testJson)
	fmt.Println("receve a message : " , testJson)
	time.Sleep(1 * time.Second)
	//处理你们的业务。。。。
	if testJson.UserId == 11 {
		//用户新增逻辑。。。。。
	}
	if testJson.UserId == 22 {
		//用户更新逻辑。。。。。
	}
	rand.Seed(time.Now().Unix())
	if rand.Intn(3) > 1 {
		var p *int
		*p = 0
		return true
	}
	return false
}
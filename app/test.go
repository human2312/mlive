package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"time"
)



func Test(c *gin.Context)  {


}

func Test2(c *gin.Context)  {
	//建立一个队列，把可以的layerupid存放
	layerSlice := []int64{21,22,23}
	totalUser := int64(100000)
	for i := int64(5);i <= totalUser;i++ {
		rand.Seed(time.Now().UnixNano())
		ranLen := rand.Intn(len(layerSlice))
		upId := layerSlice[ranLen]


		log.Println("reg....",ranLen,"value:",upId,"count",i)

		//新userId增加到layerSlice里面去
		layerSlice = append(layerSlice, i)
		//log.Println(layerSlice)
	}
	log.Println(layerSlice)
	log.Println("done....")
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": 1,
	})
}
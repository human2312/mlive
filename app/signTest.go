package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mlive/library/sign"
)

func SingTest(c *gin.Context)  {
	MapList := make(map[string]interface{})
	MapList["ddffdf"] = "dffdaxd"
	MapList["ddfdfdfe3fdf"] = "dfadf4e5"

	sign := sign.Set(MapList,"v1")
	fmt.Println("执行")
	fmt.Println(sign)

	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  sign,
		"data": "",
	})
}

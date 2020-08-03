package app

import (
	"github.com/gin-gonic/gin"
	"mlive/service/user/dao"
	"mlive/util"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": 1,
	})
}

var (
	info        dao.User
	maxChildNum int64   = 3
	ExcelUtil           = new(util.ExcelUtil)
	perIncome   float64 = 10000 //单用户收入
	income      float64 = 0
	expend      float64 = 0
	runing      bool    = false
	layerSlice          = []int64{1}
)

func Run(c *gin.Context) {

	return
}





//EchoTitle 输出title
func EchoTitle() {
	outfit := [9]string{
		"no",
		"inviter",
		"invitee",
		"profit",
		"income",
		"expend",
		"payoffs",
		"name",
		"value",
	}
	ExcelUtil.AppendFile(outfit)
}

//AppendFile 追加内容
func AppendFile(outfit [9]string) {
	ExcelUtil.AppendFile(outfit)
}


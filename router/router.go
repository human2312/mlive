package router

import (
	"bytes"
	"io/ioutil"
	"mlive/library/sign"
	"mlive/util"
	"time"

	"github.com/gin-gonic/gin"

	"mlive/app"
	"mlive/dao"
	"mlive/library/logger"
)

var (
	r *gin.Engine
)

func Init() *gin.Engine {
	dao.NewDao()
	util.NewMQ()

	r = gin.New()
	// 全局中间件 日志 和 recovery
	r.Use(logReq, gin.Recovery(), cors)

	r.GET("/ping", app.Ping)


	api := r.Group("/api")
	{
		test := api.Group("/test")
		{
			test.POST("/test", app.Test)
			test.POST("/run", app.Run)
			test.POST("/send", app.SendTest)
			test.POST("/sign", app.SingTest)
			test.POST("/admin", app.AdminTest)
		}

		shop := api.Group("/shop")
		{
			//需要验签
			shop.Use(sign.Check)
			shop.POST("/list", app.ShopList)
			shop.POST("/add", app.ShopAdd)
			shop.POST("/edit", app.ShopEdit)
			shop.POST("/del", app.ShopDel)
		}
		// 云仓
		cloud := api.Group("/cloud")
		{
			cloud.Use(sign.Check)
			cloud.POST("/list",app.CloudList)
			cloud.POST("/handle",app.HandleCloud)
			cloud.POST("/submit",app.SubmitCloud)
			cloud.POST("/recovery",app.Recovery)
			cloud.POST("/info",app.CloudInfo)
		}
		// 用户
		user := api.Group("/user")
		{
			//user.Use(sign.Check)
			user.POST("/statistics/info",app.GetUserStatisticsInfo)
			user.GET("/mq/submit",app.MqUser)

			user.POST("/invite/tree/list",sign.Check,app.GetInviteTree)
		}
	}

	//r.POST("/mq/user/receive",app.MqUserReceive)

	admin := r.Group("/admin")
	{
		test := admin.Group("/test")
		{
			test.POST("/level/set_one",app.SetCompanyAccount)
		}
		user := admin.Group("/user")
		{
			//需要验签
			user.Use(sign.Admin)
			user.GET("/account/list",app.GetUserAccountList)
			user.GET("/list",app.GetUserList)
			user.GET("/info",app.GetUserInfo)
			user.GET("/level/log",app.UserLevelLog)
			user.GET("/team/list",app.GetUserTeamList)
			user.POST("/status/update",app.SaveStatus)
			user.POST("/team/update",app.SaveUserTeam)
			user.POST("/level/update",app.SaveUserLevel)
		}
		cloud := admin.Group("/cloud")
		{
			cloud.POST("/Refund",app.Refund)
			cloud.Use(sign.Admin)
			cloud.POST("/save/number",app.SaveCloudNumber)
			cloud.GET("/pay/log",app.GetCloudPayLog)
			cloud.GET("/all/log",app.GetCloudAllLog)

		}

		// 分润配置
		profitConfig := admin.Group("/profitConfig", sign.Admin)
		{
			profitConfig.GET("/list", app.Prftcfg.List)
			profitConfig.POST("/update", app.Prftcfg.Update)
		}
	}
	return r
}

// 处理跨域请求,支持options访问
func cors(c *gin.Context) {
	method := c.Request.Method
	// 核心处理方式
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "DNT,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Admin-Token,X-MPMALL-Sign,X-MPMALL-SignVer,X-MPMALL-Token,X-MPMALL-APPVer")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.JSON(200, "")
	}
	// 处理请求
	c.Next()
}

func logReq(c *gin.Context) {
	start := time.Now()
	method := c.Request.Method
	url := c.Request.URL
	header := c.Request.Header
	body := ""
	// 把request的内容读取出来
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		body = string(bodyBytes)
	}
	// 把刚刚读出来的再写进去
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	clientip := c.ClientIP()
	c.Next()
	status := c.Writer.Status()
	errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()
	// bodySize := c.Writer.Size()
	latency := time.Now().Sub(start)

	logger.Rprintln("method:", method, "|url:", url, "|header:", header, "|body:", body, "|clientip:", clientip, "|status:", status, "|latency:", latency, "|errMsg:", errMsg)
}

package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mlive/service/admin"
	"mlive/util"
)

func AdminTest(c *gin.Context)  {
	adminToken := c.Request.Header.Get("Admin-Token")
	if adminToken == "" {
		util.Fail(c, 11111, "缺少Admin-Token参数")
		return
	}
	adminInfo ,err:= admin.CheckAdminLogin(adminToken)
	fmt.Println(adminInfo,err)
	//根据adminId 获取adminInfo
	adminInfoByid ,err2:= admin.GetAdminUserInfo(adminToken,adminInfo.Id)
	fmt.Println(adminInfo,adminInfoByid,err2)
}

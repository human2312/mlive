package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	// "log"
	"mlive/dao"
	"mlive/service/admin"
	"mlive/service/profit"
)

type profitConfig struct{}

var (
	Prftcfg = &profitConfig{}
)

type ProfitConfigResp struct {
	profit.MliveProfitConfig
	AdminUserName string `json:"adminUserName"`
}

func (pc *profitConfig) List(c *gin.Context) {
	all := (&profit.MliveProfitConfig{}).All()
	var data []ProfitConfigResp
	for _, v := range all {
		var row ProfitConfigResp
		row.MliveProfitConfig = v
		username, _ := admin.GetAdminUserInfo(c.GetHeader("Admin-Token"), v.AdminUserId)
		row.AdminUserName = username.Username
		data = append(data, row)
	}
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": gin.H{
			"items": data,
		},
	})
}

func (m ProfitConfigResp) MarshalJSON() ([]byte, error) {
	type Alias ProfitConfigResp
	return json.Marshal(&struct {
		Alias
		CreateTime string `json:"createTime"`
		UpdateTime string `json:"updateTime"`
	}{
		Alias:      Alias(m),
		CreateTime: m.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime: m.UpdateTime.Format("2006-01-02 15:04:05"),
	})
}

// type profitConfigReq struct {
// 	Id int64   `json:"id"`
// 	G1 float64 `json:"g1"`
// 	G2 float64 `json:"g2"`
// 	G3 float64 `json:"g3"`
// 	L2 float64 `json:"l2"`
// 	L3 float64 `json:"l3"`
// 	L4 float64 `json:"l4"`
// 	L5 float64 `json:"l5"`
// }

func (pc *profitConfig) Update(c *gin.Context) {
	var param profit.MliveProfitConfig
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(200, gin.H{
			"code": 80001,
			"msg":  err.Error(),
		})
		return
	}
	if param.Id <= 0 ||
		param.G1 < 0 ||
		param.G2 < 0 ||
		param.G3 < 0 ||
		param.L2 < 0 ||
		param.L3 < 0 ||
		param.L4 < 0 ||
		param.L5 < 0 {
		c.JSON(200, gin.H{
			"code": 80001,
			"msg":  "param error",
		})
		return
	}
	adminInfo, _ := admin.CheckAdminLogin(c.GetHeader("Admin-Token"))
	db := dao.Db.DB
	if db.Model(&profit.MliveProfitConfig{}).Where("id = ?", param.Id).Updates(map[string]interface{}{
		"g1":            param.G1,
		"g2":            param.G2,
		"g3":            param.G3,
		"l2":            param.L2,
		"l3":            param.L3,
		"l4":            param.L4,
		"l5":            param.L5,
		"admin_user_id": adminInfo.Id,
	}).Error != nil {
		c.JSON(200, gin.H{
			"code": 80000,
			"msg":  "fail",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
	})
}

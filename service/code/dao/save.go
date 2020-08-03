package dao

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	d "mlive/dao"
	"mlive/library/snowflakeId"
	"mlive/service/code/model"
	"time"
)

// 创建汇总-码
func (code *UserCode)CreateCodeSummary(userId int64) (int64,error) {

	var data = model.MliveUserCodeSummary{
		UserId:userId,
	}
	if err := d.Db.DB.Create(&data).Error; err !=nil{
		return 0,err
	}
	return data.Id,nil
}

//  更新汇总-码
func (code *UserCode)SaveCodeSummary(userId int64,mapData map[string]interface{})(bool,error)  {
	if err := d.Db.DB.Model(&model.MliveUserCodeSummary{}).Where("user_id=?",userId).Updates(mapData).Error; err !=nil {
		return false,err
	}else{
		return true,nil
	}
}

func (code *UserCode)SaveCodeSummaryColumn(userId int64,mapData map[string]interface{})(bool,error)  {
	//mapData["quantity"] = gorm.Expr("quantity-?",1)
	if err := d.Db.DB.Model(&model.MliveUserCodeSummary{}).Where("user_id=?",userId).UpdateColumn(mapData).Error; err !=nil {
		return false,err
	}else{
		return true,nil
	}
}


// 名额处理
// userId:用户id,ty
func (code *UserCode)CreateCodeLog(userId int64,ty int64,tyCode int64,useUserId int64,number int64,channel int64,status int64,adminUserId int64)(int64,error)  {

	no := snowflakeId.GetSnowflakeId()
	var data = model.MliveUserCodeLog{
		UserId:userId,
		No:no,
		Type:ty,
		CodeType:tyCode,
		UseUserId:useUserId,
		Number:number,
		Channel:channel,
		Status:status,
		AdminUserId:adminUserId,
		CreateTime:time.Now(),
		UpdateTime:time.Now(),
	}
	if err := d.Db.DB.Create(&data).Error; err !=nil{
		return 0,err
	}
	if  data.Id > 0 {
		//
		return no,nil
	}else{
		return 0,errors.New("create code log error")
	}
}
package dao

import (
	"errors"
	d "mlive/dao"
	"mlive/service/code/model"
)
type UserCode struct {

}

// 获取汇总信息
func (code *UserCode)GetCodeSummaryInfo(userId int64)(*model.MliveUserCodeSummary,error)  {
	if userId <= 0 {
		return &model.MliveUserCodeSummary{},errors.New("请输入用户id")
	}
	var data = model.MliveUserCodeSummary{}
	if err := d.Db.DB.First(&data).Where("user_id=?",userId).Error; err != nil{
		return &model.MliveUserCodeSummary{},err
	}
	return  &data,nil
}

// 获取总数
func (code *UserCode)GetCodeLogCount(userId int64)(int64,error)  {
	if userId <= 0 {
		return 0,errors.New("请输入用户id")
	}
	var totalScore int64 = 0
	if err := d.Db.DB.Find(&model.MliveUserCodeLog{}).Where("user_id = ?",userId).Count(&totalScore).Error; err != nil {
		return 0,err
	}
	return totalScore,nil
}

type CodeLogNum struct {
	Number int64 `json:"number"`
}

// 获取日志名额数
func (code *UserCode)GetCodeLogNum(userId int64,codeType int64)(int64,error)  {
	if userId <= 0 {
		return 0,errors.New("请输入用户id")
	}
	var Results CodeLogNum
	d.Db.DB.Table(model.MliveUserCodeLog{}.TableName()).Where("user_id = ?",userId).Where("code_type=?",codeType).Where("status=?",1).
		Select("sum(number) as number").Scan(&Results)
	return Results.Number,nil
}
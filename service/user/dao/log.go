package dao

import (
	"mlive/service/user/model"
	d "mlive/dao"
	"time"
)

// 保持用户等级日志
// type:1:用户自动升级 2:后台操作升级,
func (u *User)CreateLevelLog(ty int64,adminUserId int64,Params string) (int64,error) {

	var data = model.MliveUserLevelLog{
		Type:ty,
		Params:Params,
		AdminUserId:adminUserId,
		CreateTime:time.Now(),
		UpdateTime:time.Now(),
	}
	if err := d.Db.DB.Create(&data).Error; err !=nil{
		return 0,err
	}

	return data.Id,nil
}

/**
* 获取更新等级日志列表
 */
func (u *User)GetLevelLogList(userId int64,page int,row int,orderBy string)([]*model.MliveUserLevelLog,error)  {

	var data = []*model.MliveUserLevelLog{}
	if orderBy == ""{
		orderBy = "id desc"
	}
	user := d.Db.DB.Model(&model.MliveUserLevelLog{}).Offset((page-1)*row).Limit(row).Order(orderBy).Where("status=?",1)
	if userId > 0 {
		user = user.Where("user_id = ? ",userId)
	}
	if err := user.Find(&data).Error; err != nil{
		return data,err
	}
	return  data,nil
}

/**
* 获取更新等级日志数量
 */
func (u *User)GetLevelLogCount(userId int64)(int64,error)  {
	var total int64 = 0
	user := d.Db.DB.Model(&model.MliveUserLevelLog{}).Order("id desc")
	if userId > 0 {
		user = user.Where("user_id = ? ",userId)
	}

	if err := user.Count(&total).Error; err != nil{
		return total,err
	}
	return  total,nil
}

/**
* 修改等级日志
 */
func (u *User)SaveLevelLog(id int64,mapData map[string]interface{})(bool,error)  {
	if err := d.Db.DB.Model(&model.MliveUserLevelLog{Id:id}).Updates(mapData).Error; err !=nil {
		return false,err
	}else{
		return true,nil
	}
}


// 保持用户日志
func (u *User)CreateUserLog(userId int64,ty int64,status int64,adminUserId  int64,Params string) (int64,error) {

	var data = model.MliveUserLog{
		UserId:userId,
		Type:ty,
		Status:status,
		AdminUserId:adminUserId,
		Params:Params,
		CreateTime:time.Now(),
		UpdateTime:time.Now(),
	}
	if err := d.Db.DB.Create(&data).Error; err !=nil{
		return 0,err
	}

	return data.Id,nil
}


/**
* 修改用户日志
 */
func (u *User)SaveUserLog(id int64,mapData map[string]interface{})(bool,error)  {
	if err := d.Db.DB.Model(&model.MliveUserLog{Id:id}).Updates(mapData).Error; err !=nil {
		return false,err
	}else{
		return true,nil
	}
}
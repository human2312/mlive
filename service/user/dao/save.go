package dao

import (
	d "mlive/dao"
	"mlive/service/user/model"
)



/**
* 修改用户信息
 */
func (u *User)SaveUserInfo(userId int64,mapData map[string]interface{})(bool,error)  {
	if err := d.Db.DB.Model(&model.MliveUser{Id:userId}).Updates(mapData).Error; err !=nil {
		return false,err
	}else{
		return true,nil
	}
}



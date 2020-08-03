package dao

import (
	"errors"
	d "mlive/dao"
	"mlive/service/user/model"
)


type User struct {

}

/**
获取个人信息
 */
func (u *User)GetUserInfo(userId int64)(*model.MliveUser,error)  {
	if userId <= 0 {
		return &model.MliveUser{},errors.New("请输入正确Id")
	}
	var data = model.MliveUser{}
	if err := d.Db.DB.First(&data,userId).Error; err != nil{
		return &model.MliveUser{},err
	}
	return  &data,nil
}
// 通过code获取用户信息
func (u *User)GetUserCodeInfo(inviteCode string)(*model.MliveUser,error)  {
	if inviteCode == "" {
		return &model.MliveUser{},errors.New("请输入正确邀请码")
	}
	var data = model.MliveUser{}
	if err := d.Db.DB.Where("invite_code=?",inviteCode).First(&data).Error; err != nil{
		return &model.MliveUser{},err
	}
	return  &data,nil
}
// 通过昵称找用户id
func (u *User)GetUserNicknameInfo(nickname string)(*model.MliveUser,error)  {
	if nickname == "" {
		return &model.MliveUser{},errors.New("请输入正确用户id")
	}
	var data = model.MliveUser{}
	//Where("nickname LIKE ?","%"+nickname+"%")
	if err := d.Db.DB.Where("nickname=?",nickname).First(&data).Error; err != nil{
		return &model.MliveUser{},err
	}
	return  &data,nil
}
// 通过手机找用户信息
func (u *User)GetUserMobileInfo(mobile string)(*model.MliveUser,error)  {
	if mobile == "" {
		return &model.MliveUser{},errors.New("请输入正确手机号码")
	}
	var data = model.MliveUser{}
	if err := d.Db.DB.Where("mobile=?",mobile).First(&data).Error; err != nil{
		return &model.MliveUser{},err
	}
	return  &data,nil
}
// 通过区号+手机找用户信息
func (u *User)GetUserAreaMobileInfo(areaCode string,mobile string)(*model.MliveUser,error)  {
	if areaCode == "" {
		return &model.MliveUser{},errors.New("请输入正确区号")
	}
	if mobile == "" {
		return &model.MliveUser{},errors.New("请输入正确手机号码")
	}
	var data = model.MliveUser{}
	if err := d.Db.DB.Where("area_code=?",areaCode).Where("mobile=?",mobile).First(&data).Error; err != nil{
		return &model.MliveUser{},err
	}
	return  &data,nil
}


/**
* 获取用户列表
 */
func (u *User)GetUserList(userId int64,page int,row int,orderBy string)([]*model.MliveUser,error)  {

	var data = []*model.MliveUser{}
	if orderBy == ""{
		orderBy = "id desc"
	}
	user := d.Db.DB.Model(&model.MliveUser{}).Offset((page-1)*row).Limit(row).Order(orderBy)
	if userId > 0 {
		user = user.Where("id = ? ",userId)
	}
	if err := user.Find(&data).Error; err != nil{
		return data,err
	}
	return  data,nil
}

/**
* 获取团队列表
 */
func (u *User)GetUserTeamList(userId int64,inviteUpId int64,isTeam int64,typeLevel int,level int64)([]*model.MliveUser,error)  {

	var data = []*model.MliveUser{}
	user := d.Db.DB.Model(&model.MliveUser{}).Order("id desc")
	if userId > 0 {
		user = user.Where("id = ? ",userId)
	}
	if inviteUpId > 0 {
		user = user.Where("invite_up_id = ? ",inviteUpId)
	}
	if isTeam >= 0 {
		user = user.Where("is_team = ? ", isTeam)
	}
	if typeLevel == 1 {
		user = user.Where("level > ?  ", level)
	}else if typeLevel == 2 {
		user = user.Where("level <= ?  ", level)
	}
	if err := user.Find(&data).Error; err != nil{
		return data,err
	}
	return  data,nil
}

/**
* 获取用户总数量
 */
func (u *User)GetUserCount(userId int64,inviteUpId int64)(int64,error)  {
	var total int64 = 0
	user := d.Db.DB.Model(&model.MliveUser{}).Order("id desc")
	if userId > 0 {
		user = user.Where("id = ? ",userId)
	}
	if inviteUpId > 0 {
		user = user.Where("invite_up_id = ? ",inviteUpId)
	}
	if err := user.Count(&total).Error; err != nil{
		return total,err
	}
	return  total,nil
}

// 获取邀请列表
func (u *User)getInviteList(userVal []string)([]*model.MliveUser,error)  {

	var data = []*model.MliveUser{}
	useLen	 := len(userVal)
	if useLen > 0 {
		user :=  d.Db.DB.Model(&model.MliveUser{})
		if useLen == 1 {
			user = user.Where("invite_up_id = ?",userVal)
		}else{
			user = user.Where("invite_up_id in (?)",userVal)
		}
		if err :=  user.Find(&data).Error ; err != nil {
			return data,err
		}
	}
	return data,nil
}
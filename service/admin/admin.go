package admin

// @Time : 2020年3月9日15:09:04
// @Author : Lemyhello
// @Desc: admin

import (
	"mlive/dao"
	"mlive/library/error"
)

type Admin struct {
	Id          int
	Username    string
	LastLoginIP string
	Avatar      string
	AddTime     string
	UpdateTime  string
	Deleted     int
	Version     int
	RoleID      int
}

//CheckAdminLogin 根据登录信息拿到当前用户
func CheckAdminLogin(adminToken string) (Admin, *error.Err) {
	admininfo := Admin{}
	err := &error.Err{}
	adminRec ,err:= dao.GetAdminInfo(adminToken)
	if err == nil {
		if adminRec.Code == 10000 && adminRec.Data.ID > 0 {
			admininfo.Id = adminRec.Data.ID
			admininfo.Username = adminRec.Data.Username
			err = nil
		} else {
			err = error.New(adminRec.Code,adminRec.Msg)
		}
	}
	return admininfo,err
}

//GetAdminUserInfo 根据id拿到管理用户信息
func GetAdminUserInfo(adminToken string, id int) (admininfo Admin ,err *error.Err) {
	adminRec ,neterr:= dao.GetAdminInfoById(adminToken, id)
	if neterr == nil {
		if adminRec.Code == 10000 && adminRec.Data.ID > 0 {
			admininfo.Id = adminRec.Data.ID
			admininfo.Username = adminRec.Data.Username
			admininfo.LastLoginIP = adminRec.Data.LastLoginIP
			admininfo.Avatar = adminRec.Data.Avatar
			admininfo.AddTime = adminRec.Data.AddTime
			admininfo.UpdateTime = adminRec.Data.UpdateTime
			admininfo.Deleted = adminRec.Data.Deleted
			admininfo.Version = adminRec.Data.Version
			admininfo.RoleID = adminRec.Data.RoleID
			err = nil
		} else {
			err = error.New(adminRec.Code,adminRec.Msg)
		}
	} else {
		err = neterr
	}
	return admininfo,err
}

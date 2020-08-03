package profit

import (
	// "log"
	"mlive/library/logger"
	"mlive/service/user/arango"
	"mlive/service/user/dao"
)

type user struct {
	id         int64
	level      int
	inviteUpId int64
}

var argoUser = &arango.ArangoUser{}
var daoUser = &dao.User{}

func updateLevel(id int64, suitLevel int, orderNo string) bool {
	if id <= 0 {
		return false
	}
	success, err := daoUser.SaveLevel(id, int64(suitLevel), 2, int64(suitLevel), 0, "", orderNo)
	logger.Iprintln(logFlag, " savelevel ", id, suitLevel, orderNo, success, err)
	return success
}

// 获取用户信息
func getUser(id int64) (u user) {
	if id <= 0 {
		return u
	}
	au, err := argoUser.Info(id)
	logger.Iprintln(logFlag, " getUser ", id, au, err)
	return infoReset(au)
}

// 获取一代用户
func getLeft1User(id int64) (u user) {
	if id <= 0 {
		return u
	}
	au, err := argoUser.GetSameLevelReward(id)
	logger.Iprintln(logFlag, " getLeft1User ", id, au, err)
	return infoReset(*au)
}

// // 获取越级用户
// func getLeapFrogUser(id int64) (u user) {
// 	if id <= 0 {
// 		return u
// 	}
// 	au, err := argoUser.GetLeapFrogInfo(id)
// 	logger.Iprintln(logFlag, " getLeapFrogUser ", id, au, err)
// 	return infoReset(*au)
// }

// 获取越级用户
func getLeapFrogUser(u user) (leapFrogUser user) {
	if u.id <= 0 {
		return
	}
	if u.inviteUpId <= 0 {
		return
	}
	inviteUp1User := getUser(u.inviteUpId)

	if inviteUp1User.id <= 0 || inviteUp1User.level == 0 || inviteUp1User.level >= u.level {
		return
	}

	logger.Iprintln(logFlag, " getLeapFrogUser ", u, inviteUp1User)

	return inviteUp1User
}

// 获取分润上级用户
func getUp1User(id int64) (u user) {
	if id <= 0 {
		return u
	}
	au, err := argoUser.GetMoneySuperior(id)
	logger.Iprintln(logFlag, " getUp1User ", id, au, err)
	return infoReset(*au)
}

func infoReset(a arango.MyUserInfo) (u user) {
	if a.IsCompany == 1 || a.Id == 0 {
		u.id = 0
		u.level = 0
		u.inviteUpId = 0
		return
	}
	u.id = a.Id
	u.level = int(a.Level)
	u.inviteUpId = a.InviteUpId
	return
}

// mlive/service/user/dao/save.go
// 等级修改:SaveLevel
// mlive/service/user/arango/info.go
// 越级上级信息:GetLeapFrogInfo

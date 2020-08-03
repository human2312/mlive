/**
 * 面膜99 分润规则
 */
package profit

import (
	// "log"
	"mlive/library/logger"
	cloud "mlive/service/cloud/dao"
)

// profitType
// 差价 190
// 全部返还 180
// 拉新给予 171  +10
// 拉新扣减 172  -10
// 进货奖励 160
// 联创1 151
// 联创2 152
// 运费扣减 141

type maskAmountCfg struct {
	buyPrice float64
}

type maskAllLevelCfg map[int]maskAmountCfg

var (
	recommendNew float64 = 10
	l5LeftAmount float64 = 0.5
)

// num 件数
func (m maskAllLevelCfg) cal(first bool, orderNo string, num int, id int64) bool {
	if id <= 0 || num <= 0 {
		return true
	}
	u := getUser(id)
	if u.id <= 0 {
		return true
	}
	var allPl []MliveProfitLog
	up1User := getUp1User(id)
	inviteUp1User := getUser(u.inviteUpId)

	if up1User.id > 0 {
		cac, err := getCarriageConfig()
		if err != nil {
			return false
		}
		allPl = append(allPl, cac.calCarriage(orderNo, up1User, num, "mask"))
	}

	if u.level == 0 && first {
		if inviteUp1User.id > 0 && inviteUp1User.level == 0 {
			allPl = append(allPl, m.calRecommendNewIn(orderNo, inviteUp1User))
			if up1User.id > 0 {
				allPl = append(allPl, m.calRecommendNewOut(orderNo, up1User))
			}
		}
	}

	if inviteUp1User.id > 0 && inviteUp1User.level > u.level {
		useNum := m.useCode(orderNo, 3*num, u, inviteUp1User)
		if useNum > 0 {
			allPl = append(allPl, m.calBack(orderNo, u.level, int(useNum/3), inviteUp1User))
			if inviteUp1User.level >= 3 && inviteUp1User.level <= 5 {
				allPl = append(allPl, getSaleAmountConfig().calStockSale(orderNo, inviteUp1User, useNum, "mask"))
			}
			num = num - int(useNum/3)
		}
	}

	if u.level == 0 && num > 0 && inviteUp1User.level == 0 && inviteUp1User.id > 0 {
		inviteUp2User := getUser(inviteUp1User.inviteUpId)
		if inviteUp2User.id > 0 && inviteUp2User.level >= 2 {
			useNum := m.useCode(orderNo, 3*num, u, inviteUp2User)
			if useNum > 0 {
				allPl = append(allPl, m.calBack(orderNo, u.level, int(useNum/3), inviteUp2User))
				if inviteUp2User.level >= 3 && inviteUp2User.level <= 5 {
					allPl = append(allPl, getSaleAmountConfig().calStockSale(orderNo, inviteUp2User, useNum, "mask"))
				}
				num = num - int(useNum/3)
			}
		}
	}

	if num > 0 {
		allPl = append(allPl, m.calAllDiff(orderNo, num, u, up1User)...)
	}

	validPlList := getValid(allPl)

	logger.Iprintln(logFlag, orderNo, " profit result ", validPlList)

	if insertList(orderNo, validPlList) {
		return true
	} else {
		logger.Eprintln(logFlag, orderNo, " profit result in mysql failed ", validPlList)
		return false
	}
}

// )SaveReduceCloudStoragGoodsNums(userId int64,num int64,orderNo string)(int64,error)
func (m maskAllLevelCfg) useCode(orderNo string, num int, from, to user) int {
	if from.id <= 0 || to.id <= 0 || from.level >= to.level {
		return 0
	}
	userNum, _ := (&cloud.CloudStorage{}).SaveReduceCloudStoragGoodsNums(to.id, from.id, int64(num), orderNo)
	if userNum < 0 {
		userNum = -userNum
	}
	return int(userNum)
}

func (m maskAllLevelCfg) calRecommendNewIn(orderNo string, to user) (pl MliveProfitLog) {
	if to.id <= 0 || to.level != 0 {
		return pl
	}
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 171,
		Amount:     recommendNew,
	}
}

func (m maskAllLevelCfg) calRecommendNewOut(orderNo string, to user) (pl MliveProfitLog) {
	if to.id <= 0 || to.level <= 0 {
		return pl
	}
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 172,
		Amount:     -recommendNew,
	}
}

func (m maskAllLevelCfg) calBack(orderNo string, level int, num int, to user) (pl MliveProfitLog) {
	if to.id <= 0 || num <= 0 && level < 0 || level > 5 {
		return pl
	}
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 180,
		Amount:     m[level].buyPrice * float64(num),
	}
}

func (m maskAllLevelCfg) calOneDiff(orderNo string, num int, from, to user) (pl MliveProfitLog) {
	if from.id <= 0 || to.id <= 0 || num <= 0 {
		return pl
	}
	amount := (m[from.level].buyPrice - m[to.level].buyPrice) * float64(num)
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 190,
		Amount:     amount,
	}
}

func (m maskAllLevelCfg) calAllDiff(orderNo string, num int, u, firstUp1User user) (allPl []MliveProfitLog) {
	if num <= 0 || u.id <= 0 {
		return
	}
	from := u
	to := firstUp1User
	for {
		if to.id <= 0 || to.level <= from.level {
			break
		}
		pl := m.calOneDiff(orderNo, num, from, to)
		allPl = append(allPl, pl)
		if to.level == 5 && pl.Amount > 0 {
			left1User := getLeft1User(to.id)
			if left1User.id > 0 && left1User.level == 5 {
				allPl = append(allPl, m.calLeft1(orderNo, num, left1User))
				left2User := getLeft1User(left1User.id)
				if left2User.id > 0 && left2User.level == 5 {
					allPl = append(allPl, m.calLeft2(orderNo, num, left2User))
				}
			}
		}
		from = to
		to = getUp1User(from.id)
	}
	return
}

func (m maskAllLevelCfg) calLeft1(orderNo string, num int, u user) (pl MliveProfitLog) {
	if num <= 0 || u.id <= 0 || u.level != 5 {
		return
	}
	amount := l5LeftAmount * 3 * float64(num)
	return MliveProfitLog{
		UserId:     u.id,
		OrderNo:    orderNo,
		ProfitType: 151,
		Amount:     amount,
	}
}

func (m maskAllLevelCfg) calLeft2(orderNo string, num int, u user) (pl MliveProfitLog) {
	if num <= 0 || u.id <= 0 || u.level != 5 {
		return
	}
	amount := l5LeftAmount * 3 * float64(num)
	return MliveProfitLog{
		UserId:     u.id,
		OrderNo:    orderNo,
		ProfitType: 152,
		Amount:     amount,
	}
}

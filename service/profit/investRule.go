package profit

import (
	// "log"
	"mlive/library/logger"
	cloud "mlive/service/cloud/dao"
)

// profitType
// 差价 290
// 全部返还 280
// 越级奖励 270
// 一代 211
// 二代 212
// 三代 213
// 进货奖励 260
// 联创1 251
// 联创2 252

// if cloudType == 0  {
//   num = 3
// }else if cloudType == 1 {
//   num = 21
// }else if cloudType == 2 { //店长
//   num = 29
// }else if cloudType == 3 {
//   num = 117
// }else if cloudType == 4 { //1764
//   num = 294
// }else if cloudType == 5 { //5292
//   num = 17647
// }

var levelRelateNum map[int]int = map[int]int{
	0: 3,
	1: 21,
	2: 29,
	3: 117,
	4: 294,
	5: 17647,
}

type investAmountCfg struct {
	buyPrice    float64
	generation1 float64
	generation2 float64
	generation3 float64
	diff        map[int]float64
}

type investAllLevelCfg map[int]investAmountCfg

var (
	generationTotal int = 3
)

func (i investAllLevelCfg) cal(orderNo string, id int64, targetLevel int) bool {
	generationNum := generationTotal
	if id <= 0 || targetLevel <= 0 {
		return true
	}
	u := getUser(id)
	if u.id <= 0 {
		return true
	}
	var allPl []MliveProfitLog
	if u.level < targetLevel {
		if !updateLevel(id, targetLevel, orderNo) {
			logger.Eprintln(logFlag, " orderNo ", orderNo, " userId ", id, " updateLevel falied ")
			return false
		}
		u = getUser(id)
		if u.id <= 0 || u.level < 1 {
			return true
		}

		useCode := false

		inviteUp1User := getUser(u.inviteUpId)
		if inviteUp1User.id > 0 && inviteUp1User.level > u.level {
			useNum := i.useCode(orderNo, targetLevel, u, inviteUp1User)
			if useNum > 0 {
				useCode = true
				// 全部返还
				allPl = append(allPl, i.calBack(orderNo, targetLevel, inviteUp1User))
				if inviteUp1User.level >= 3 && inviteUp1User.level <= 5 {
					allPl = append(allPl, getSaleAmountConfig().calStockSale(orderNo, inviteUp1User, useNum, "invest"))
				}
			} else {
				leapFrogUser := getLeapFrogUser(u)
				if leapFrogUser.id > 0 {
					// 越级奖励
					allPl = append(allPl, i.calLeapFrog(orderNo, targetLevel, leapFrogUser))
					// 差价 ================
					allPl = append(allPl, i.calAllGenerationAndAllDiff(false, orderNo, generationNum, u, targetLevel, true)...)
				} else {
					// 差价 + 代奖励 =========
					allPl = append(allPl, i.calAllGenerationAndAllDiff(false, orderNo, generationNum, u, targetLevel, false)...)
				}
			}

		} else {
			leapFrogUser := getLeapFrogUser(u)
			if leapFrogUser.id > 0 {
				// 越级奖励
				allPl = append(allPl, i.calLeapFrog(orderNo, targetLevel, leapFrogUser))
				// 差价 ================
				allPl = append(allPl, i.calAllGenerationAndAllDiff(false, orderNo, generationNum, u, targetLevel, true)...)
			} else {
				// 差价 + 代奖励 =========
				allPl = append(allPl, i.calAllGenerationAndAllDiff(false, orderNo, generationNum, u, targetLevel, false)...)
			}
		}

		if useCode {
			(&cloud.CloudStorage{}).SaveReduceCloudStoragNums(u.id, 1, int64(targetLevel), inviteUp1User.id, orderNo, 2, 0)
		} else {
			(&cloud.CloudStorage{}).SaveReduceCloudStoragNums(u.id, 1, int64(targetLevel), 0, orderNo, 2, 0)
		}

	} else {

		useCode := false

		inviteUp1User := getUser(u.inviteUpId)
		if inviteUp1User.id > 0 && inviteUp1User.level > u.level {
			useNum := i.useCode(orderNo, targetLevel, u, inviteUp1User)
			if useNum > 0 {
				useCode = true
				// 全部返还
				allPl = append(allPl, i.calBack(orderNo, targetLevel, inviteUp1User))
				if inviteUp1User.level >= 3 && inviteUp1User.level <= 5 {
					allPl = append(allPl, getSaleAmountConfig().calStockSale(orderNo, inviteUp1User, useNum, "invest"))
				}
			} else {
				// 复购 差价 + 代奖励
				allPl = append(allPl, i.calAllGenerationAndAllDiff(false, orderNo, generationNum, u, targetLevel, false)...)
			}
		} else {
			allPl = append(allPl, i.calAllGenerationAndAllDiff(false, orderNo, generationNum, u, targetLevel, false)...)
		}

		if useCode {
			(&cloud.CloudStorage{}).SaveReduceCloudStoragNums(u.id, 1, int64(targetLevel), inviteUp1User.id, orderNo, 2, 0)
		} else {
			(&cloud.CloudStorage{}).SaveReduceCloudStoragNums(u.id, 1, int64(targetLevel), 0, orderNo, 2, 0)
		}

	}

	// 记录入库

	validPlList := getValid(allPl)

	logger.Iprintln(logFlag, orderNo, " profit result ", validPlList)

	if insertList(orderNo, validPlList) {
		return true
	} else {
		logger.Eprintln(logFlag, orderNo, " profit result in mysql failed ", validPlList)
		return false
	}

}

func (i investAllLevelCfg) calAllGenerationAndAllDiff(buyDown bool, orderNo string, generationNum int, u user, targetLevel int, half bool) (allPl []MliveProfitLog) {
	if u.id <= 0 {
		return
	}
	var to user
	if !buyDown {
		to = getLeft1User(u.id)
		for {
			if generationNum <= 0 {
				break
			}
			if to.id <= 0 {
				break
			}
			switch generationNum {
			case 3:
				allPl = append(allPl, i.calGeneration1(orderNo, targetLevel, to, half))
			case 2:
				allPl = append(allPl, i.calGeneration2(orderNo, targetLevel, to))
			case 1:
				allPl = append(allPl, i.calGeneration3(orderNo, targetLevel, to))
			}
			generationNum--
			to = getLeft1User(to.id)
		}
	}

	from := u
	to = getUp1User(u.id)
	for {
		if to.id <= 0 {
			break
		}
		if generationNum > 0 {
			switch generationNum {
			case 3:
				allPl = append(allPl, i.calGeneration1(orderNo, targetLevel, to, half))
			case 2:
				allPl = append(allPl, i.calGeneration2(orderNo, targetLevel, to))
			case 1:
				allPl = append(allPl, i.calGeneration3(orderNo, targetLevel, to))
			}
			generationNum--
		}

		pl := i.calOneDiff(orderNo, targetLevel, from, to)
		allPl = append(allPl, pl)

		if num, ok := levelRelateNum[targetLevel]; ok {
			if to.level == 5 && pl.Amount > 0 {
				left1User := getLeft1User(to.id)
				if left1User.id > 0 && left1User.level == 5 {
					allPl = append(allPl, i.calLeft1(orderNo, num, left1User))
					left2User := getLeft1User(left1User.id)
					if left2User.id > 0 && left2User.level == 5 {
						allPl = append(allPl, i.calLeft2(orderNo, num, left2User))
					}
				}
			}
		}

		from = to
		to = getUp1User(from.id)
	}
	return
}

func (i investAllLevelCfg) calGeneration1(orderNo string, targetLevel int, to user, half bool) (pl MliveProfitLog) {
	if targetLevel < 1 || targetLevel > 5 || to.id <= 0 {
		return pl
	}
	gc := i[targetLevel]
	var amount float64
	if half {
		amount = gc.generation1 / 2
	} else {
		amount = gc.generation1
	}
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 211,
		Amount:     amount,
	}
}

func (i investAllLevelCfg) calGeneration2(orderNo string, targetLevel int, to user) (pl MliveProfitLog) {
	if targetLevel < 1 || targetLevel > 5 || to.id <= 0 {
		return pl
	}
	gc := i[targetLevel]
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 212,
		Amount:     gc.generation2,
	}
}

func (i investAllLevelCfg) calGeneration3(orderNo string, targetLevel int, to user) (pl MliveProfitLog) {
	if targetLevel < 1 || targetLevel > 5 || to.id <= 0 {
		return pl
	}
	gc := i[targetLevel]
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 213,
		Amount:     gc.generation3,
	}
}

func (i investAllLevelCfg) calLeapFrog(orderNo string, targetLevel int, to user) (pl MliveProfitLog) {
	if targetLevel < 1 || targetLevel > 5 || to.id <= 0 {
		return pl
	}
	amount := i[targetLevel].generation1 / 2
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 270,
		Amount:     amount,
	}

}

func (i investAllLevelCfg) useCode(orderNo string, targetLevel int, from, to user) int {
	if from.id <= 0 || to.id <= 0 || from.level >= to.level {
		return 0
	}
	num, _ := (&cloud.CloudStorage{}).SaveReduceCloudStoragNums(to.id, 2, int64(targetLevel), from.id, orderNo, 2, 0)
	if num < 0 {
		num = -num
	}
	return int(num)
}

func (i investAllLevelCfg) calBack(orderNo string, targetLevel int, to user) (pl MliveProfitLog) {
	if to.id <= 0 || targetLevel < 1 || targetLevel > 5 {
		return pl
	}
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 280,
		Amount:     i[targetLevel].buyPrice,
	}
}

func (i investAllLevelCfg) calOneDiff(orderNo string, targetLevel int, from, to user) (pl MliveProfitLog) {
	if from.id <= 0 || to.id <= 0 || from.level >= to.level {
		return pl
	}
	dc := i[targetLevel].diff
	amount := dc[to.level] - dc[from.level]
	return MliveProfitLog{
		UserId:     to.id,
		OrderNo:    orderNo,
		ProfitType: 290,
		Amount:     amount,
	}
}

func (i investAllLevelCfg) calLeft1(orderNo string, num int, u user) (pl MliveProfitLog) {
	if num <= 0 || u.id <= 0 || u.level != 5 {
		return
	}
	amount := l5LeftAmount * float64(num)
	return MliveProfitLog{
		UserId:     u.id,
		OrderNo:    orderNo,
		ProfitType: 251,
		Amount:     amount,
	}
}

func (i investAllLevelCfg) calLeft2(orderNo string, num int, u user) (pl MliveProfitLog) {
	if num <= 0 || u.id <= 0 || u.level != 5 {
		return
	}
	amount := l5LeftAmount * float64(num)
	return MliveProfitLog{
		UserId:     u.id,
		OrderNo:    orderNo,
		ProfitType: 252,
		Amount:     amount,
	}
}

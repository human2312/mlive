package profit

import (
	// "log"
	"mlive/library/logger"
)

// profitType
// 面膜销售 160
// 套餐 260

type stockSaleAmountConfig struct {
	stockSale float64
}

type stockSaleAllLevelConfig map[int]stockSaleAmountConfig

func (ssalc stockSaleAllLevelConfig) calStockSale(orderNo string, u user, num int, saleType string) (pl MliveProfitLog) {
	if len(ssalc) <= 0 {
		return
	}
	if u.id <= 0 || u.level < 3 || u.level > 5 {
		return
	}
	if num == 0 {
		return
	} else if num < 0 {
		num = -num
	}
	profitType := 0
	switch saleType {
	case "mask":
		profitType = 160
	case "invest":
		profitType = 260
	}
	amount := ssalc[u.level].stockSale * float64(num)
	return MliveProfitLog{
		UserId:     u.id,
		OrderNo:    orderNo,
		ProfitType: profitType,
		Amount:     amount,
	}
}

func (ssalc stockSaleAllLevelConfig) cal(orderNo string, id int64, useNum int, saleType string) bool {
	if id <= 0 && useNum <= 0 {
		return true
	}
	u := getUser(id)
	if u.id <= 0 {
		return true
	}
	var allPl []MliveProfitLog
	if u.level >= 3 && u.level <= 5 {
		allPl = append(allPl, ssalc.calStockSale(orderNo, u, useNum, saleType))
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

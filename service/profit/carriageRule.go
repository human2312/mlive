package profit

import (
// "log"
// "mlive/library/logger"
)

// profitType
// 面膜销售 141

type carriageAmountConfig struct {
	price float64
}

// num 件数
func (cac carriageAmountConfig) calCarriage(orderNo string, u user, num int, saleType string) (pl MliveProfitLog) {
	if cac.price <= 0 || u.id <= 0 {
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
		profitType = 141
	case "invest":
		// profitType = 260
	}
	amount := cac.price * float64(num)
	return MliveProfitLog{
		UserId:     u.id,
		OrderNo:    orderNo,
		ProfitType: profitType,
		Amount:     -amount,
	}
}

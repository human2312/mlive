package profit

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	// "log"
	"mlive/library/logger"
	"mlive/library/sign"
)

var (
	ver string = "v1"
)

type DepositOrder struct {
	DepositNo string `json:"depositNo"`
	UserId    int64  `json:"userId"`
	BuyType   int    `json:"buyType"`
	CouponSn  string `json:"couponSn"`
}

type CodeMsg struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type DepositOrderRes struct {
	CodeMsg
	DepositOrder `json:"data"`
}

type GoodsOrder struct {
	OrderInfo struct {
		OrderNo      string `json:"orderNo"`
		UserId       int64  `json:"userId"`
		BuyTypeCount int    `json:"buyTypeCount"`
		BuyType      int    `json:"buyType"`
		CouponSn     string `json:"couponSn"`
	} `json:"orderInfo"`
	OrderGoods []struct {
		Number  int   `json:"number"`
		GoodsId int64 `json:"goodsId"`
		// Price  float64 `json:"price"`
	} `json:"orderGoods"`
}

type GoodsOrderRes struct {
	CodeMsg
	GoodsOrder `json:"data"`
}

// type MaskBuyTotal struct {
// 	BuyTypeCount int `json:"buyTypeCount"`
// }

// type MaskBuyTotalRes struct {
// 	CodeMsg
// 	MaskBuyTotal `json:"data"`
// }

// func maskBuyTotal(id int64) (total int) {
// 	url := viper.GetString("java.java_maskBuyTotal_url")
// 	params := map[string]interface{}{
// 		"userId": id,
// 	}
// 	sign := sign.Set(params, ver)
// 	res, err := resty.New().
// 		NewRequest().
// 		SetHeader("X-MPMALL-SignVer", ver).
// 		SetHeader("X-MPMALL-Sign", sign).
// 		SetBody(params).
// 		Post(url)

// 	if err != nil {
// 		logger.Eprintln(" maskBuyTotal res: ", id, res, err)
// 	} else {
// 		logger.Iprintln(" maskBuyTotal res: ", id, res, err)

// 		var mbtr MaskBuyTotalRes
// 		json.Unmarshal(res.Body(), &mbtr)
// 		if mbtr.Code == 10000 {
// 			return mbtr.MaskBuyTotal.BuyTypeCount
// 		}
// 	}
// 	return 0
// }

func goodsOrder(orderNo string) (g GoodsOrder) {
	url := viper.GetString("java.java_goodsOrderDetail_url")
	params := map[string]interface{}{
		"orderNo": orderNo,
	}
	sign := sign.Set(params, ver)
	res, err := resty.New().
		NewRequest().
		SetHeader("X-MPMALL-SignVer", ver).
		SetHeader("X-MPMALL-Sign", sign).
		SetBody(params).
		Post(url)

	if err != nil {
		logger.Eprintln(logFlag, " orderNo res: ", orderNo, res, err)
	} else {
		logger.Iprintln(logFlag, " orderNo res: ", orderNo, res, err)

		var gor GoodsOrderRes
		err := json.Unmarshal(res.Body(), &gor)
		if err != nil {
			logger.Eprintln(logFlag, " orderNo Unmarshal res: ", orderNo, res, err)
			return
		}
		if gor.Code == 10000 {
			return gor.GoodsOrder
		}
	}
	return g
}

func depositOrder(orderNo string) (do DepositOrder) {
	url := viper.GetString("java.java_depositOrderDetail_url")
	params := map[string]interface{}{
		"depositNo": orderNo,
	}
	sign := sign.Set(params, ver)
	res, err := resty.New().
		NewRequest().
		SetHeader("X-MPMALL-SignVer", ver).
		SetHeader("X-MPMALL-Sign", sign).
		SetBody(params).
		Post(url)

	if err != nil {
		logger.Eprintln(logFlag, " orderNo res: ", orderNo, res, err)
	} else {
		logger.Iprintln(logFlag, " orderNo res: ", orderNo, res, err)

		var dor DepositOrderRes
		err := json.Unmarshal(res.Body(), &dor)
		if err != nil {
			logger.Eprintln(logFlag, " orderNo Unmarshal res: ", orderNo, res, err)
			return
		}
		if dor.Code == 10000 {
			return dor.DepositOrder
		}
	}
	return do
}

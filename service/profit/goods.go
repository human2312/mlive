/**
 * 商品详情
 */
package profit

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	// "log"
	"mlive/library/logger"
	// "mlive/library/sign"
	"strconv"
)

type PhpCodeMsg struct {
	Msg  string `json:"msg"`
	Code string `json:"code"`
}

type GoodsDetailPhp struct {
	Id       string `json:"id"`
	Plain    string `json:"plain"`
	Plus     string `json:"plus"`
	Train    string `json:"train"`
	Serv     string `json:"serv"`
	Partner  string `json:"partner"`
	Director string `json:"director"`
}

type GoodsDetail struct {
	Id       int64   `json:"id"`
	Plain    float64 `json:"plain"`
	Plus     float64 `json:"plus"`
	Train    float64 `json:"train"`
	Serv     float64 `json:"serv"`
	Partner  float64 `json:"partner"`
	Director float64 `json:"director"`
}

type GoodsDetailPhpRes struct {
	PhpCodeMsg
	Data GoodsDetailPhp `json:"Data"`
}

func goodsDetail(id int64) (gd GoodsDetail) {
	url := viper.GetString("php.php_goodsDetail_url")
	params := map[string]string{
		"id": strconv.FormatInt(id, 10),
	}
	res, err := resty.New().
		NewRequest().
		// SetHeader("X-MPMALL-SignVer", ver).
		// SetHeader("X-MPMALL-Sign", sign).
		SetQueryParams(params).
		Get(url)

	if err != nil {
		logger.Eprintln(logFlag, " goodsDetail id res: ", id, res, err)
	} else {
		logger.Iprintln(logFlag, " goodsDetail id res: ", id, res, err)

		var gdr GoodsDetailPhpRes
		err := json.Unmarshal(res.Body(), &gdr)
		if err != nil {
			logger.Eprintln(logFlag, " goodsDetail Unmarshal res: ", id, res, err)
			return
		}
		if gdr.Code == "10000" {
			gd = goodsDetailPhp2GoodsDetail(gdr.Data)
		}
	}
	return gd
}

func goodsDetailPhp2GoodsDetail(gh GoodsDetailPhp) (g GoodsDetail) {
	g.Id, _ = strconv.ParseInt(gh.Id, 10, 64)
	g.Plain, _ = strconv.ParseFloat(gh.Plain, 64)
	g.Plus, _ = strconv.ParseFloat(gh.Plus, 64)
	g.Train, _ = strconv.ParseFloat(gh.Train, 64)
	g.Serv, _ = strconv.ParseFloat(gh.Serv, 64)
	g.Partner, _ = strconv.ParseFloat(gh.Partner, 64)
	g.Director, _ = strconv.ParseFloat(gh.Director, 64)
	return
}

type InvestGoodsDetailPhp struct {
	Id       string `json:"id"`
	Plain    string `json:"plain"`
	Plus     string `json:"plus"`
	Train    string `json:"train"`
	Serv     string `json:"serv"`
	Partner  string `json:"partner"`
	Director string `json:"director"`
}

type InvestGoodsDetailPhpRes struct {
	PhpCodeMsg
	Data struct {
		Items []InvestGoodsDetailPhp `json:"items"`
	} `json:"data"`
}

type InvestGoodsDetail struct {
	Id       int64   `json:"id"`
	Plain    float64 `json:"plain"`
	Plus     float64 `json:"plus"`
	Train    float64 `json:"train"`
	Serv     float64 `json:"serv"`
	Partner  float64 `json:"partner"`
	Director float64 `json:"director"`
}

func investGoodsDetail() (igd InvestGoodsDetail) {
	url := viper.GetString("php.php_investGoodsDetail_url")

	res, err := resty.New().
		NewRequest().
		Get(url)

	if err != nil {
		logger.Eprintln(logFlag, " investGoodsDetail res: ", res, err)
	} else {
		logger.Iprintln(logFlag, " investGoodsDetail res: ", res, err)

		var igdr InvestGoodsDetailPhpRes
		err := json.Unmarshal(res.Body(), &igdr)
		if err != nil {
			logger.Eprintln(logFlag, " investGoodsDetail Unmarshal res: ", res, err)
			return
		}
		if igdr.Code == "10000" && len(igdr.Data.Items) > 0 {
			igd = investGoodsDetailPhp2investGoodsDetail(igdr.Data.Items[0])
		}
	}
	return
}

func investGoodsDetailPhp2investGoodsDetail(gh InvestGoodsDetailPhp) (g InvestGoodsDetail) {
	g.Id, _ = strconv.ParseInt(gh.Id, 10, 64)
	g.Plain, _ = strconv.ParseFloat(gh.Plain, 64)
	g.Plus, _ = strconv.ParseFloat(gh.Plus, 64)
	g.Train, _ = strconv.ParseFloat(gh.Train, 64)
	g.Serv, _ = strconv.ParseFloat(gh.Serv, 64)
	g.Partner, _ = strconv.ParseFloat(gh.Partner, 64)
	g.Director, _ = strconv.ParseFloat(gh.Director, 64)
	return
}

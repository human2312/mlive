package profit

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	// "log"
	"mlive/library/logger"
	// "mlive/library/sign"
)

type CarriageResp struct {
	Carriage `json:"data"`
	CodeMsg
}

type Carriage struct {
	Id    int     `json:"id"`
	Price float64 `json:"price"`
}

func carriageDetail() (c Carriage) {
	url := viper.GetString("java.java_carriage_url")

	res, err := resty.New().
		NewRequest().
		Get(url)

	if err != nil {
		logger.Eprintln(logFlag, " carriage detail ", res, err)
	} else {
		logger.Iprintln(logFlag, " carriage detail ", res, err)

		var cr CarriageResp
		err := json.Unmarshal(res.Body(), &cr)
		if err != nil {
			logger.Eprintln(logFlag, " carriage Unmarshal res: ", res, err)
			return
		}
		if cr.Code == 10000 {
			return cr.Carriage
		}
	}
	return
}

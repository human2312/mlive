package profit

import (
	"testing"
	// "time"

	"mlive/dao"
	"mlive/util"

	"mlive/library/config"
	"mlive/library/logger"

	"github.com/spf13/viper"
	// "strconv"

	"github.com/go-resty/resty/v2"
)

func init() {
	viper.AddConfigPath("../../")
	config.Init()
	logger.Init()
	dao.NewDao()
	util.NewMQ()

}

func Test_resty(t *testing.T) {
	client := resty.New()

	url := viper.GetString("java.java_depositOrderDetail_url")
	params := map[string]interface{}{
		"depositNo": "2019111500013872",
	}

	res, err := client.R().SetHeader("X-MPMALL-SignVer", "v1").SetBody(params).Post(url)

	t.Log(res, err)
}

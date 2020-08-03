package logger

import (
	"testing"

	"mlive/library/config"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("../../")
	config.Init()
	Init()
}

func Test_log(t *testing.T) {
	for i := 0; i < 10000; i++ {
		Iprintln("a log test ", i)
	}
}

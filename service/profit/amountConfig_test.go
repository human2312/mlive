package profit

import (
	"testing"
)

func Test_getMaskAmountConfig(t *testing.T) {
	t.Log(getMaskAmountConfig(1))
}

func Test_getInvestAmountConfig(t *testing.T) {
	t.Log(getInvestAmountConfig())
}

func Test_all(t *testing.T) {
	t.Log((&MliveProfitConfig{}).All())
}

func Test_getCarriageConfig(t *testing.T) {
	t.Log(getCarriageConfig())
}

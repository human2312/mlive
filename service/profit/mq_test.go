package profit

import (
	"testing"
)

func Test_OrderReceive(t *testing.T) {
	OrderReceive()
}

// func Test_publishProfitQueue(t *testing.T) {
// 	publishProfitQueue("no002")
// }

func Test_handle(t *testing.T) {
	strBtye := []byte(`
{"orderType":3,"orderNo":"2020040200227425","uuid":"d4bbe014-7482-49f5-8e76-e44bf4349da5"}

`)

	t.Log(handle(strBtye))

}

func Test_publishProfitQueue(t *testing.T) {
	publishProfitQueue("2020032600225904")
}

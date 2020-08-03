package profit

import (
	"testing"
)

func Test_calCarriage(t *testing.T) {
	cac, err := getCarriageConfig()
	t.Log(cac, err)
	t.Log(cac.calCarriage("----test---", user{1, 1, 1}, 10, "mask"))
}

package profit

import (
	"testing"
)

func Test_maskRuleCal(t *testing.T) {
	malc.cal(true, "10000000", 2, 11)
}

func Test_calRecommendNewIn(t *testing.T) {
	t.Log(malc.calRecommendNewIn("test001", user{
		id:         11,
		level:      0,
		inviteUpId: 1,
	}))
}

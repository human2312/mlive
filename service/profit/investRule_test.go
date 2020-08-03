package profit

import (
	"testing"
)

func Test_investRuleCal(t *testing.T) {
	iac := getInvestAmountConfig()
	t.Log(iac.calAllGenerationAndAllDiff(false, "lf000001", 3, user{1215, 1, 0}, 1, false))
}

package profit

import (
	// "encoding/json"
	"errors"
	"mlive/dao"
	"time"
)

var (
	logFlag string = "<order2profit>"
)

func getCarriageConfig() (cac carriageAmountConfig, err error) {
	c := carriageDetail()
	if c.Id == 0 {
		return cac, errors.New(" get carriage error ")
	}
	return carriageAmountConfig{
		price: c.Price,
	}, nil
}

func getSaleAmountConfig() (ssalc stockSaleAllLevelConfig) {
	ssalc = stockSaleAllLevelConfig{
		3: stockSaleAmountConfig{
			stockSale: 2,
		},
		4: stockSaleAmountConfig{
			stockSale: 4,
		},
		5: stockSaleAmountConfig{
			stockSale: 5,
		},
	}
	return
}

func getMaskAmountConfig(id int64) (mlc maskAllLevelCfg) {
	gd := goodsDetail(id)
	if gd.Id <= 0 {
		return
	}

	mlc = make(maskAllLevelCfg)

	mlc[0] = maskAmountCfg{
		buyPrice: gd.Plain * 3,
	}
	mlc[1] = maskAmountCfg{
		buyPrice: gd.Plus * 3,
	}
	mlc[2] = maskAmountCfg{
		buyPrice: gd.Train * 3,
	}
	mlc[3] = maskAmountCfg{
		buyPrice: gd.Serv * 3,
	}
	mlc[4] = maskAmountCfg{
		buyPrice: gd.Partner * 3,
	}
	mlc[5] = maskAmountCfg{
		buyPrice: gd.Director * 3,
	}

	return

}

func getInvestAmountConfig() (iac investAllLevelCfg) {
	igd := investGoodsDetail()
	if igd.Id <= 0 {
		return
	}
	allPc := (&MliveProfitConfig{}).All()
	if len(allPc) != 5 {
		return
	}

	iac = investAllLevelCfg{
		1: investAmountCfg{
			buyPrice:    igd.Plus,
			generation1: allPc[0].G1,
			generation2: allPc[0].G2,
			generation3: allPc[0].G3,

			diff: map[int]float64{
				1: 0,
				2: 0 + allPc[0].L2,
				3: 0 + allPc[0].L2 + allPc[0].L3,
				4: 0 + allPc[0].L2 + allPc[0].L3 + allPc[0].L4,
				5: 0 + allPc[0].L2 + allPc[0].L3 + allPc[0].L4 + allPc[0].L5,
			},
		},
		2: investAmountCfg{
			buyPrice:    igd.Train,
			generation1: allPc[1].G1,
			generation2: allPc[1].G2,
			generation3: allPc[1].G3,

			diff: map[int]float64{
				1: 0,
				2: 0 + allPc[1].L2,
				3: 0 + allPc[1].L2 + allPc[1].L3,
				4: 0 + allPc[1].L2 + allPc[1].L3 + allPc[1].L4,
				5: 0 + allPc[1].L2 + allPc[1].L3 + allPc[1].L4 + allPc[1].L5,
			},
		},
		3: investAmountCfg{
			buyPrice:    igd.Serv,
			generation1: allPc[2].G1,
			generation2: allPc[2].G2,
			generation3: allPc[2].G3,

			diff: map[int]float64{
				1: 0,
				2: 0 + allPc[2].L2,
				3: 0 + allPc[2].L2 + allPc[2].L3,
				4: 0 + allPc[2].L2 + allPc[2].L3 + allPc[2].L4,
				5: 0 + allPc[2].L2 + allPc[2].L3 + allPc[2].L4 + allPc[2].L5,
			},
		},
		4: investAmountCfg{
			buyPrice:    igd.Partner,
			generation1: allPc[3].G1,
			generation2: allPc[3].G2,
			generation3: allPc[3].G3,

			diff: map[int]float64{
				1: 0,
				2: 0 + allPc[3].L2,
				3: 0 + allPc[3].L2 + allPc[3].L3,
				4: 0 + allPc[3].L2 + allPc[3].L3 + allPc[3].L4,
				5: 0 + allPc[3].L2 + allPc[3].L3 + allPc[3].L4 + allPc[3].L5,
			},
		},
		5: investAmountCfg{
			buyPrice:    igd.Director,
			generation1: allPc[4].G1,
			generation2: allPc[4].G2,
			generation3: allPc[4].G3,

			diff: map[int]float64{
				1: 0,
				2: 0 + allPc[4].L2,
				3: 0 + allPc[4].L2 + allPc[4].L3,
				4: 0 + allPc[4].L2 + allPc[4].L3 + allPc[4].L4,
				5: 0 + allPc[4].L2 + allPc[4].L3 + allPc[4].L4 + allPc[4].L5,
			},
		},
	}

	return
}

func (m *MliveProfitConfig) All() (allPc []MliveProfitConfig) {
	db := dao.Db.DB
	db.Order("goods_level asc").Find(&allPc)
	return
}

type MliveProfitConfig struct {
	Id          int64     `json:"id"`
	GoodsLevel  int       `json:"goodsLevel"`
	G1          float64   `json:"g1"`
	G2          float64   `json:"g2"`
	G3          float64   `json:"g3"`
	L2          float64   `json:"l2"`
	L3          float64   `json:"l3"`
	L4          float64   `json:"l4"`
	L5          float64   `json:"l5"`
	AdminUserId int       `json:"-"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

package profit

import (
	"encoding/json"
	"github.com/spf13/viper"
	// "log"
	"mlive/library/logger"
	cloud "mlive/service/cloud/dao"
	"mlive/util"
	"strconv"
	"time"
)

var (
	RabbitMQ               = new(util.RabbitMQ)
	orderMqName            string
	profitMqName           string
	monitorUserLevelMqName string
)

type OrderMsg struct {
	OrderNo   string `json:"orderNo"`
	OrderType int    `json:"orderType"`
	// UpUserId  int64  `json:"upUserId"`
}

// 订单类型
// -1:商品99购买
// 1:VIP499
// 2:赚播店长999
// 3:赚播总监9999
// 4:赚播合伙人30000
// 5:赚播联创89100
// {"orderType":-1,"orderNo":"2020031400224476","uuid":"0ec0f952-4f4f-4c90-8c61-b80223783770"}

//读到一条消息
func OrderReceive() {
	orderMqName = viper.GetString("rabbitmqueue.order")

	RabbitMQ.Consume(orderMqName, handle)
}

//
// 线下的 加加这两个
// 用户帮下级升级vip/店长/总监/合伙人 调用下面的
// cloud.SaveCloudLogByOrder(0,云仓编号, map[string]interface{}{
//   "order_no":订单号,
//   "update_time":time.Now(),
// })

// 用户赠送商品,调用
// SaveGiveGoodsStatus(cloudNo int64)
//
func handle(msg []byte) bool {
	logger.Iprintln(logFlag, " mq ", orderMqName, " receive a msg: ", string(msg))

	var om OrderMsg
	json.Unmarshal(msg, &om)

	orderNo := om.OrderNo
	orderType := om.OrderType

	if orderNo == "" {
		logger.Iprintln(logFlag, " mq ", orderMqName, " receive a invalid msg: ", string(msg))
		return true
	}

	done := false

	od := &MliveOrderDone{
		OrderNo: orderNo,
	}

	if od.have() {
		logger.Iprintln(logFlag, " mq ", orderMqName, " msg has been done: ", string(msg))
		return true
	}

	switch {
	case orderType == -1:
		g := goodsOrder(orderNo)

		// 只有goodsId == 1 的商品才参与分润
		if len(g.OrderGoods) == 1 {
			orderGoods := g.OrderGoods[0]
			goodsId := orderGoods.GoodsId
			if goodsId != 1 {
				logger.Iprintln(logFlag, " mq ", orderMqName, " receive a not profit msg: ", string(msg))
				return true
			}
		}

		if g.OrderInfo.BuyType == 2 {
			// 进货奖励
			ssalc := getSaleAmountConfig()
			couponSn := g.OrderInfo.CouponSn
			if couponSn == "" {
				logger.Iprintln(logFlag, " mq ", orderMqName, " receive a invalid msg: ", string(msg))
				return true
			}
			csno, err := strconv.ParseInt(couponSn, 10, 64)
			if csno > 0 && err == nil {
				cs, err := (&cloud.CloudStorage{}).GetCloudStoragOrderInfo(0, csno)
				if err != nil {
					logger.Iprintln(logFlag, " mq ", orderMqName, " receive a invalid msg: ", string(msg), " couponSn res err: ", err)
				} else {
					// 通知计光
					(&cloud.CloudStorage{}).SaveGiveGoodsStatus(csno)
					userId := cs.UserId
					useNum := cs.Number
					if useNum < 0 {
						useNum = -useNum
					}
					done = ssalc.cal(orderNo, userId, int(useNum), "mask")
				}
			}
		} else {
			id := g.OrderInfo.UserId
			if id > 0 {
				total := g.OrderInfo.BuyTypeCount
				if total > 0 {
					first := false
					if total == 1 {
						first = true
					}
					if len(g.OrderGoods) == 1 {
						orderGoods := g.OrderGoods[0]
						num := orderGoods.Number
						goodsId := orderGoods.GoodsId
						if num > 0 && goodsId > 0 {
							mac := getMaskAmountConfig(goodsId)
							if len(mac) > 0 {
								done = mac.cal(first, orderNo, int(num/3), id)
							}
						}
					}

				}
			}
		}
	case orderType >= 1 && orderType <= 5:
		d := depositOrder(orderNo)
		if d.BuyType == 2 {
			// 进货奖励
			ssalc := getSaleAmountConfig()
			couponSn := d.CouponSn
			if couponSn == "" {
				logger.Iprintln(logFlag, " mq ", orderMqName, " receive a invalid msg: ", string(msg))
				return true
			}
			csno, err := strconv.ParseInt(couponSn, 10, 64)
			if csno > 0 && err == nil {
				cs, err := (&cloud.CloudStorage{}).GetCloudStoragOrderInfo(0, csno)
				if err != nil {
					logger.Iprintln(logFlag, " mq ", orderMqName, " receive a invalid msg: ", string(msg), " couponSn res err: ", err)
				} else {
					// 通知计光
					(&cloud.CloudStorage{}).SaveCloudLogByOrder(0, csno, map[string]interface{}{
						"order_no":    orderNo,
						"update_time": time.Now(),
					})

					userId := cs.UserId
					useNum := cs.Number
					if useNum < 0 {
						useNum = -useNum
					}
					done = ssalc.cal(orderNo, userId, int(useNum), "invest")
				}
			}
		} else {
			id := d.UserId
			if id > 0 {
				targetLevel := orderType
				iac := getInvestAmountConfig()
				if len(iac) > 0 {
					done = iac.cal(orderNo, id, targetLevel)
				}
			}
		}
	default:
		logger.Iprintln(logFlag, " mq ", orderMqName, " receive a invalid msg: ", string(msg))
		return true
	}

	if done {
		publishProfitQueue(orderNo)
		logger.Iprintln(logFlag, " mq orderNo:  ", orderNo, " deal success ", string(msg))
		return true
	}
	logger.Eprintln(logFlag, " mq orderNo:  ", orderNo, " deal fail ", string(msg))

	return false
}

type Item struct {
	OrderNo      string  `json:"orderNo"`
	ProfitUserId int64   `json:"profitUserId"`
	ProfitAmount float64 `json:"profitAmount"`
	ProfitType   int     `json:"profitType"`
	RewardType   int     `json:"rewardType"`
	OperateType  int     `json:"operateType"`
}

type ProfitMsg struct {
	Items []Item `json:"items"`
}

func publishProfitQueue(orderNo string) {
	profitMqName = viper.GetString("rabbitmqueue.result")
	profitList := ListByOrderNo(orderNo)
	if len(profitList) == 0 {
		return
	}
	var back ProfitMsg
	back.Items = rebuildProfitList(profitList)
	if len(back.Items) == 0 {
		return
	}
	backByte, _ := json.Marshal(back)
	err := RabbitMQ.Publish(profitMqName, string(backByte))
	if err == nil {
		logger.Iprintln(logFlag, profitMqName+" send a msg success: "+string(backByte))
	} else {
		logger.Eprintln(logFlag, profitMqName+" send a msg failed: "+string(backByte))
	}

}

func rebuildProfitList(in []MliveProfitLog) (out []Item) {
	for _, v := range in {
		out = append(out, rebuildProfitRow(v))
	}
	return
}

func rebuildProfitRow(in MliveProfitLog) (out Item) {
	out.OrderNo = in.OrderNo
	out.ProfitUserId = in.UserId
	out.ProfitAmount = in.Amount
	out.OperateType = 1
	if out.ProfitAmount < 0 {
		out.ProfitAmount = -out.ProfitAmount
		out.OperateType = 2
	}
	switch in.ProfitType {
	case 190:
		out.ProfitType = 4
		out.RewardType = 1
	case 180:
		out.ProfitType = 8
		out.RewardType = 1
	case 171:
		out.ProfitType = 6
		out.RewardType = 1
	case 172:
		out.ProfitType = 7
		out.RewardType = 1
	case 160:
		out.ProfitType = 9
		out.RewardType = 1
	case 151:
		out.ProfitType = 10
		out.RewardType = 1
	case 152:
		out.ProfitType = 10
		out.RewardType = 1
	case 141:
		out.ProfitType = 11
		out.RewardType = 1

	case 290:
		out.ProfitType = 4
		out.RewardType = 2
	case 280:
		out.ProfitType = 8
		out.RewardType = 2
	case 270:
		out.ProfitType = 5
		out.RewardType = 2
	case 211:
		out.ProfitType = 1
		out.RewardType = 2
	case 212:
		out.ProfitType = 2
		out.RewardType = 2
	case 213:
		out.ProfitType = 3
		out.RewardType = 2
	case 260:
		out.ProfitType = 9
		out.RewardType = 2
	case 251:
		out.ProfitType = 10
		out.RewardType = 2
	case 252:
		out.ProfitType = 10
		out.RewardType = 2

	}

	return
}

// 收益类型
// 1:一代
// 2:二代
// 3:三代
// 4:差价
// 5:越级
// 6:拉新给予（+10）
// 7:拉新扣减（-10）
// 8:消耗码的全部返还
// 9:进货（身份）奖励
//
// 奖励类型 1:商品销售奖励 2:服务商销售奖励
//
// 操作类型 1：增加收益 2：减少收益
//
//
//
// 差价 190
// 全部返还 180
// 拉新给予 171  +10
// 拉新扣减 172  -10
// 进货奖励 160
//
//
// 差价 290
// 全部返还 280
// 越级奖励 270
// 一代 211
// 二代 212
// 三代 213
// 进货奖励 260

// type UserLevelMsg struct {
// 	UserId    int64  `json:"userId"`
// 	Type      int    `json:"type"`
// 	No        string `json:"no"`
// 	OrderType string `json:"orderType"`
// }

// func publishUserLevelQueue(orderType string, orderNo string, userId int64) {
// 	monitorUserLevelMqName = viper.GetString("rabbitmqueue.monitorUserLevel")
// 	ulmsg := UserLevelMsg{
// 		UserId:    userId,
// 		Type:      2,
// 		No:        orderNo,
// 		OrderType: orderType,
// 	}
// 	backByte, _ := json.Marshal(ulmsg)
// 	err := RabbitMQ.Publish(monitorUserLevelMqName, string(backByte))
// 	if err == nil {
// 		logger.Iprintln(monitorUserLevelMqName + " send a msg success: " + string(backByte))
// 	} else {
// 		logger.Eprintln(monitorUserLevelMqName + " send a msg failed: " + string(backByte))
// 	}
// }

package profit

import (
	"mlive/dao"
	"mlive/library/logger"
	"time"
)

type MliveProfitLog struct {
	Id         int64
	UserId     int64
	ProfitType int
	Amount     float64
	OrderNo    string
	CreateTime time.Time
	UpdateTime time.Time
}

type MliveOrderDone struct {
	Id         int64
	OrderNo    string
	CreateTime time.Time
	UpdateTime time.Time
}

func (od *MliveOrderDone) have() bool {
	db := dao.Db.DB
	var o MliveOrderDone
	if db.Where(od).First(&o).RecordNotFound() {
		return false
	}
	return true
}

func (od *MliveOrderDone) Insert() bool {
	db := dao.Db.DB
	od.CreateTime = time.Now()
	od.UpdateTime = time.Now()
	if err := db.Create(od).Error; err != nil {
		return false
	}
	return true
}

func (pl *MliveProfitLog) Insert() bool {
	db := dao.Db.DB
	pl.CreateTime = time.Now()
	pl.UpdateTime = time.Now()
	if err := db.Create(pl).Error; err != nil {
		return false
	}
	return true
}

func (pl *MliveProfitLog) List() []MliveProfitLog {
	var list []MliveProfitLog
	db := dao.Db.DB
	db.Find(&list, pl)
	return list
}

func ListByOrderNo(orderNo string) (list []MliveProfitLog) {
	if orderNo == "" {
		return list
	}
	pl := &MliveProfitLog{
		OrderNo: orderNo,
	}
	list = pl.List()
	return list
}

func insertList(orderNo string, list []MliveProfitLog) bool {
	db := dao.Db.DB
	tx := db.Begin()

	od := MliveOrderDone{
		OrderNo:    orderNo,
		UpdateTime: time.Now(),
		CreateTime: time.Now(),
	}

	if tx.Create(&od).Error != nil {
		tx.Rollback()
		logger.Iprintln(logFlag, orderNo, " has been done ")
		return false
	}

	if len(list) > 0 {
		for _, v := range list {
			v.CreateTime = time.Now()
			v.UpdateTime = time.Now()
			if tx.Create(&v).Error != nil {
				tx.Rollback()
				return false
			}
		}
	}
	tx.Commit()
	return true
}

func getValid(list []MliveProfitLog) (validList []MliveProfitLog) {
	if len(list) == 0 {
		return
	}
	for _, v := range list {
		if v.UserId > 0 && v.Amount != 0 {
			validList = append(validList, v)
		}
	}
	return
}

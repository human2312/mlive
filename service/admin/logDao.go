package admin

import (
	"mlive/dao"
	"time"
)

/**
 * @Author: Lemyhello
 * @Description:
 * @File:  logDao
 * @Version: X.X.X
 * @Date: 2020/4/3 上午11:20
 */

var (
	table = "mlive_rabbitmq_log"
)

type RabbitmqLog struct {
	Id         int64
	Queue      string
	Uid        string
	Types      string
	Msg        string
	Result     int
	CreateTime time.Time
	UpdateTime time.Time
}

//Add 增加
func (s *RabbitmqLog) Add(rabbitmqLog RabbitmqLog) bool {
	var data = RabbitmqLog{
		Queue:      rabbitmqLog.Queue,
		Uid:        rabbitmqLog.Uid,
		Types:      rabbitmqLog.Types,
		Msg:        rabbitmqLog.Msg,
		Result:     rabbitmqLog.Result,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if err := dao.Db.DB.Table(table).Create(&data).Error; err != nil {
		return false
	}
	return true
}

//Updates 更新
func (s *RabbitmqLog) Updates(rabbitmqLog RabbitmqLog) bool {
	if err := dao.Db.DB.Table(table).Where("uid = ?" , rabbitmqLog.Uid).Updates(map[string]interface{}{"result" : rabbitmqLog.Result,"update_time" : time.Now()}).Error; err !=nil {
		return false
	}
	return true
}
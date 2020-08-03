package shop

import (
	"mlive/dao"
	"time"
)

var (
	table = "mlive_shop"
)

type Shop struct {
	Id         int64
	UserId     int64
	Title      string
	Subtitle   string
	Page       string
	Subpage    string
	Price      float64
	Status     int
	CreateTime time.Time
	UpdateTime time.Time
}

func (s *Shop) List() []Shop {
	db    := dao.Db.DB
	var shop []Shop
	db.Table(table).Find(&shop)
	return shop
}

func (s *Shop) GetOne() Shop  {
	db    := dao.Db.DB
	var shop Shop
	db.Table(table).First(&shop,"{\"id\": \"1\"}")
	return shop
}

func (s *Shop) Add() bool {
	db    := dao.Db.DB
	s.Status = 0
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	if err := db.Table(table).Create(s).Error; err != nil {
		return false
	}
	return true
}

func (s *Shop) Edit() bool {
	db    := dao.Db.DB
	s.UpdateTime = time.Now()
	if err := db.Table(table).Model(s).Update().Error; err != nil {
		return false
	}
	return true
}

func (s *Shop) Del(id int64) bool {
	db    := dao.Db.DB
	s.Id = id
	if err := db.Table(table).Delete(s).Error; err != nil {
		return false
	}
	return true
}



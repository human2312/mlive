package mysql

import (
	"github.com/jinzhu/gorm"
	"log"
	"strings"
	"time"
	//"context"
)



type DbConf struct {
	DSN string
	Active int
	Idle int
	IdleTimeout int
}



func NewMysql(c DbConf) (db *gorm.DB,e error) {

	DSN := c.DSN
	Active := c.Active
	Idle := c.Idle
	IdleTimeout := c.IdleTimeout
	db, err := gorm.Open("mysql", DSN)
	if err != nil {
		return nil,err
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(Idle)
	db.DB().SetMaxOpenConns(Active)
	db.DB().SetConnMaxLifetime(time.Duration(IdleTimeout) / time.Second)
	//db.SetLogger(ormLog{})
	return db,nil
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Printf(strings.Repeat("%v ", len(v)), v...)
}


//func InitORM(d *gorm.DB) {
//
//	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
//		return defaultTableName
//	}
//	d.LogMode(false)
//	return
//}


//func  Ping(c context.Context,d *gorm.DB) (err error) {
//	if d.DB() != nil {
//		err = d.DB().PingContext(c)
//	}
//	return
//}
//
//func Close(d *gorm.DB)  {
//	if d.DB() != nil {
//		d.DB().Close()
//	}
//}
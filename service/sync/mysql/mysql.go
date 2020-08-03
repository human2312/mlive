package mysql

import (
	"context"
	"github.com/jinzhu/gorm"
	"time"
)

// 迁移幸福引力用户，所有用户初始等级为粉丝

type Dao struct {
	DB 		 *gorm.DB
}

type DbConf struct {
	DSN string
	Active int
	Idle int
	IdleTimeout int
}


func New(c DbConf)(d *Dao,e error){

	newDb,err := NewMysql(c)
	if err != nil {
		return nil,err
	}
	d = &Dao{
		DB:newDb,
	}
	d.initORM()
	return d,nil
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

func (d *Dao) initORM() {

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}
	d.DB.LogMode(false)
	return
}


func (d *Dao) Ping(c context.Context) (err error) {
	if d.DB != nil {
		err = d.DB.DB().PingContext(c)
	}
	return
}

func (d *Dao) Close()  {
	if d.DB != nil {
		d.DB.Close()
	}
}
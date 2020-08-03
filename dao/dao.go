package dao

import (
	"context"
	"github.com/arangodb/go-driver"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type Dao struct {
	DB       *gorm.DB
	ArangoDb driver.Database
	Redis    *redis.Pool
}

func New() (d *Dao) {

	d = &Dao{
		DB:       NewMysql(),
		ArangoDb: NewArangoDb(),
		Redis:    NewRedis(),
	}
	d.initORM()
	return
}

func (d *Dao) Ping(c context.Context) (err error) {
	if d.DB != nil {
		err = d.DB.DB().PingContext(c)
	}
	return
}

func (d *Dao) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}

package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	mysql "mlive/library/mysql"
	"strconv"
)

func NewMysql() (db *gorm.DB) {
	var conf mysql.DbConf
	DSN 			:= viper.GetString("mysql.dsn")
	Active, _ 		:= strconv.Atoi(viper.GetString("mysql.active"))
	Idle, _ 		:= strconv.Atoi(viper.GetString("mysql.idle"))
	IdleTimeout, _  := strconv.Atoi(viper.GetString("mysql.idleTimeout"))

	conf.DSN = DSN
	conf.Active = Active
	conf.Idle = Idle
	conf.IdleTimeout = IdleTimeout
	db,_ = mysql.NewMysql(conf)
	return
}
func (d *Dao) initORM() {

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}
	d.DB.LogMode(false)
	return
}

//func NewMysql() (db *gorm.DB) {
//
//	DSN := viper.GetString("mysql.dsn")
//	Active, _ := strconv.Atoi(viper.GetString("mysql.active"))
//	Idle, _ := strconv.Atoi(viper.GetString("mysql.idle"))
//	IdleTimeout, _ := strconv.Atoi(viper.GetString("mysql.idleTimeout"))
//	db, err := gorm.Open("mysql", DSN)
//	if err != nil {
//		logger.Eprintf("db dsn(%s) error:(%v) ", DSN, err)
//	}
//	db.SingularTable(true)
//	db.DB().SetMaxIdleConns(Idle)
//	db.DB().SetMaxOpenConns(Active)
//	db.DB().SetConnMaxLifetime(time.Duration(IdleTimeout) / time.Second)
//	//db.SetLogger(ormLog{})
//	return
//}
//

//type ormLog struct{}
//
//func (l ormLog) Print(v ...interface{}) {
//	log.Printf(strings.Repeat("%v ", len(v)), v...)
//}
